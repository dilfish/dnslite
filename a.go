// sean 2021 at shanghai

package dnslite

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

type AHandler struct{}

func (a *AHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.A, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeA
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].A = net.ParseIP(record.A)
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (a *AHandler) CheckRecord(record *DNSRecord) error {
	ip := net.ParseIP(record.A)
	if ip == nil {
		log.Println("check ipv4 error:", record.A)
		return ErrBadValue
	}
	ip = ip.To4()
	if ip == nil {
		log.Println("not ipv4")
		return ErrBadValue
	}
	return nil
}
