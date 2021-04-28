package dnslite

import (
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
		rr[idx].Hdr.Ttl = record.TTL
		rr[idx].Target = record.CNAME
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (c *CNAMEHandler) CheckRecord(record *DNSRecord) error {
	return nil
}
