// sean at shanghai
// 2021

package main

import (
	"github.com/miekg/dns"
)

type MxHandler struct{}

func (mx *MxHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.MX, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeMX
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Mx = record.MxMx
		rr[idx].Preference = record.MxPreference
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (mx *MxHandler) CheckRecord(record *DNSRecord) error {
	if record.MxPreference == 0 || record.MxMx == "" {
		return ErrBadValue
	}
	return nil
}

func (mx *MxHandler) RRToRecord(msg dns.RR) DNSRecord {
	var record DNSRecord
	v := msg.(*dns.MX)
	record.MxMx = v.Mx
	record.MxPreference = v.Preference
	record.Ttl = v.Hdr.Ttl
	return record
}
