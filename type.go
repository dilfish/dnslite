package dnslite

import "github.com/miekg/dns"

type TypeHandler interface {
	FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg
	CheckRecord(record *DNSRecord) error
}

var TypeHandlerList = map[uint16]TypeHandler{
	dns.TypeNS: &NSHandler{},
}

func CommonCheck(r *DNSRecord) error {
	if r.Name == "" {
		return ErrBadName
	}
	if r.Name[len(r.Name)-1] != '.' {
		r.Name = r.Name + "."
	}
	if r.TTL > 600 || r.TTL < 1 {
		return ErrBadTTL
	}
	return nil
}
