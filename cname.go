package main

import (
	"log"

	"github.com/miekg/dns"
)

type CNAMEHandler struct{}

func (c *CNAMEHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.CNAME, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeCNAME
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Target = record.Cname
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (c *CNAMEHandler) CheckRecord(record *DNSRecord) error {
	is := GoodName(record.Cname)
	if !is {
		log.Println("bad cname:", record.Cname)
		return ErrBadValue
	}
	record.Cname = AppendDot(record.Cname)
	return nil
}

func (c *CNAMEHandler) RRToRecord(msg dns.RR) DNSRecord {
	var record DNSRecord
	v := msg.(*dns.CNAME)
	record.Cname = v.Target
	record.Ttl = v.Hdr.Ttl
	return record
}
