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
	M *RecordManager
}

func NewHandler(conf *MongoClientConfig) *Handler {
	m := NewRecordManager(conf)
	return &Handler{M: m}
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Print("Client:", w.RemoteAddr())
	log.Println("au:", r.Authoritative, "tr:", r.Truncated, "rd:", r.RecursionDesired, "ra:", r.RecursionAvailable, "ad:", r.AuthenticatedData, "cd:", r.CheckingDisabled)
	err := ParseReqInfo(r)
	if err != nil {
		log.Println("bad dns request:", err)
		return
	}
	msg, err := h.GetRecord(r)
	if err != nil {
		log.Println("get record error:", err)
		return
	}
	w.WriteMsg(msg)
}

func (h *Handler) GetRecord(req *dns.Msg) (*dns.Msg, error) {
	name := req.Question[0].Name
	tp := dns.Type(req.Question[0].Qtype)
	log.Println("name and type:", name, tp)
	records, err := h.M.FindRecords(ReadRecordArgs{Name: name, Type: uint16(tp)})
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
