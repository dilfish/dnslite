// sean at shanghai 2021

package dnslite

import (
	"strings"

	"github.com/miekg/dns"
)

type TXTHandler struct{}

func (txt *TXTHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.TXT, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeA
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Txt = strings.Split(record.Txt, "\"")
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (txt *TXTHandler) CheckRecord(record *DNSRecord) error {
	if len(record.Txt) == 0 || len(record.Txt) > 2048 {
		return ErrBadValue
	}
	return nil
}
