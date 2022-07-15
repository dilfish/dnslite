// Copyright 2018 Sean.ZH

package main

import (
	"errors"
	"log"

	"github.com/miekg/dns"
)

func ParseReqInfo(r *dns.Msg) (err error) {
	if len(r.Question) != 1 {
		err = ErrBadQCount
		log.Println("Questions are not 1:", len(r.Question))
		return
	}
	PrintExtra(r.Extra)
	return
}

type Handler struct {
	M           DataManagerI
	ServFailMsg *dns.Msg
}

func NewHandler(conf *Config, db DataManagerI) *Handler {
	var h Handler
	h.ServFailMsg = new(dns.Msg)
	h.M = db
	return &h
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Print("Client:", w.RemoteAddr())
	log.Println("au:", r.Authoritative, "tr:", r.Truncated, "rd:", r.RecursionDesired, "ra:", r.RecursionAvailable, "ad:", r.AuthenticatedData, "cd:", r.CheckingDisabled)
	h.ServFailMsg.SetReply(r)
	err := ParseReqInfo(r)
	if err != nil {
		h.ServFailMsg.MsgHdr.Rcode = dns.RcodeBadName
		log.Println("bad dns request:", err)
		w.WriteMsg(h.ServFailMsg)
		return
	}
	rcode, msg, err := h.GetRecord(r)
	if err != nil {
		h.ServFailMsg.MsgHdr.Rcode = dns.RcodeNameError
		if err == ErrRCode {
			h.ServFailMsg.MsgHdr.Rcode = rcode
		}
		log.Println("get record error:", err)
		w.WriteMsg(h.ServFailMsg)
		return
	}
	err = w.WriteMsg(msg)
	if err != nil {
		log.Println("write msg error:", msg, err)
	}
}

func (h *Handler) GetRecord(req *dns.Msg) (int, *dns.Msg, error) {
	name := req.Question[0].Name
	rcode := IsRcode(name)
	if rcode > 0 {
		log.Println("name rcode:", name, rcode)
		return rcode, nil, ErrRCode
	}
	tp := dns.Type(req.Question[0].Qtype)
	log.Println("name and type:", name, tp)
	// find cname first
	records, err := h.M.Find(name, dns.TypeCNAME)
	if err == nil && len(records) > 0 {
		log.Println("cname found, return")
		msg := TypeHandlerList[dns.TypeCNAME].FillRecords(req, records)
		return 0, msg, nil
	}
	// then ns
	records, err = h.M.Find(name, dns.TypeNS)
	if err == nil && len(records) > 0 {
		log.Println("ns found, return")
		msg := TypeHandlerList[dns.TypeNS].FillRecords(req, records)
		return 0, msg, nil
	}
	records, err = h.M.Find(name, uint16(tp))
	if err != nil {
		log.Println("find records error:", name, tp, err)
		return 0, nil, err
	}
	if len(records) == 0 {
		if *FlagNoneProxy {
			log.Println("no record without proxy")
			return 0, nil, errors.New("no such record")
		}
		log.Println("proxy to real dns")
		usingTls := IfProxyTls(name, tp)
		return GetDataFromRealDNS(req, usingTls)
	}
	log.Println("find from cache", records)
	msg := TypeHandlerList[req.Question[0].Qtype].FillRecords(req, records)
	return 0, msg, nil
}
