// sean at shanghai 2021

package main

import (
	"github.com/miekg/dns"
)

type CAAHandler struct{}

func (cca *CAAHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.CAA, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeCAA
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Value = record.CAAValue
		rr[idx].Tag = record.CAATag
		rr[idx].Flag = record.CAAFlag
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (caa *CAAHandler) CheckRecord(record *DNSRecord) error {
	return nil
}
