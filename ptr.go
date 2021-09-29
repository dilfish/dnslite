package main

import (
	"github.com/miekg/dns"
)

type PtrHandler struct{}

func (p *PtrHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.PTR, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypePTR
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Ptr = record.Ptr
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (p *PtrHandler) CheckRecord(record *DNSRecord) error {
	if record.Ptr == "" {
		return ErrBadValue
	}
	return nil
}

func (p *PtrHandler) RRToRecord(msg dns.RR) DNSRecord {
	var record DNSRecord
	v := msg.(*dns.PTR)
	record.Ptr = v.Ptr
	record.Ttl = v.Hdr.Ttl
	return record
}
