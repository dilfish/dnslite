// Copyright 2018 Sean.ZH

package dnslite

import (
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
	M           *RecordManager
	ServFailMsg *dns.Msg
}

func NewHandler(conf *MongoClientConfig) *Handler {
	m := NewRecordManager(conf)
	return &Handler{M: m, ServFailMsg: new(dns.Msg)}
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
	msg, err := h.GetRecord(r)
	if err != nil {
		h.ServFailMsg.MsgHdr.Rcode = dns.RcodeServerFailure
		log.Println("get record error:", err)
		w.WriteMsg(h.ServFailMsg)
		return
	}
	err = w.WriteMsg(msg)
	if err != nil {
		log.Println("write msg error:", msg, err)
	}
}

func (h *Handler) GetRecord(req *dns.Msg) (*dns.Msg, error) {
	name := req.Question[0].Name
	tp := dns.Type(req.Question[0].Qtype)
	log.Println("name and type:", name, tp)
	// find cname first
	records, err := h.M.FindRecords(ReadRecordArgs{Name: name, Type: dns.TypeCNAME})
	if err == nil && len(records) > 0 {
		log.Println("cname found, return")
		msg := TypeHandlerList[dns.TypeCNAME].FillRecords(req, records)
		return msg, nil
	}
	// then ns
	records, err = h.M.FindRecords(ReadRecordArgs{Name: name, Type: dns.TypeNS})
	if err == nil && len(records) > 0 {
		log.Println("ns found, return")
		msg := TypeHandlerList[dns.TypeNS].FillRecords(req, records)
		return msg, nil
	}
	records, err = h.M.FindRecords(ReadRecordArgs{Name: name, Type: uint16(tp)})
	if err != nil {
		log.Println("find records error:", name, tp, err)
		return nil, err
	}
	if len(records) == 0 {
		log.Println("proxy to real dns")
		return GetDataFromRealDNS(req)
	}
	log.Println("find from cache")
	msg := TypeHandlerList[req.Question[0].Qtype].FillRecords(req, records)
	return msg, nil
}
