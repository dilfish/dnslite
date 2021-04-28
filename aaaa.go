package dnslite

import (
	"net"

	"github.com/miekg/dns"
)

type AAAAHandler struct{}

func (aaaa *AAAAHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.AAAA, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeAAAA
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.TTL
		rr[idx].AAAA = net.ParseIP(record.AAAA)
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (aaaa *AAAAHandler) CheckRecord(record *DNSRecord) error {
	ip := net.ParseIP(record.AAAA)
	if ip == nil {
		return ErrBadValue
	}
	ip = ip.To4()
	if ip != nil {
		return ErrBadValue
	}
	return nil
}
