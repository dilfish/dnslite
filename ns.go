package dnslite

import "github.com/miekg/dns"

type NSHandler struct{}

func (ns *NSHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	m.Authoritative = true
	rr := make([]dns.NS, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeNS
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.TTL
		rr[idx].Ns = record.NS
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (ns *NSHandler) CheckRecord(record *DNSRecord) error {
	return nil
}
