package main

import (
	"log"
	"strings"

	"github.com/miekg/dns"
)

type TypeHandler interface {
	FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg
	CheckRecord(record *DNSRecord) error
	RRToRecord(rr dns.RR) DNSRecord
}

var TypeHandlerList = map[uint16]TypeHandler{
	dns.TypeA:     &AHandler{},
	dns.TypeAAAA:  &AAAAHandler{},
	dns.TypeNS:    &NSHandler{},
	dns.TypeCNAME: &CNAMEHandler{},
	dns.TypeTXT:   &TXTHandler{},
	dns.TypeCAA:   &CAAHandler{},
	dns.TypeSVCB:  &SVCBHandler{},
	dns.TypeSOA:   &SoaHandler{},
}

func CommonCheck(r *DNSRecord) error {
	if r.Name == "" {
		log.Println("bad record name:", r.Name)
		return ErrBadName
	}
	if r.Name[len(r.Name)-1] != '.' {
		r.Name = r.Name + "."
	}
	if r.Ttl > 600 || r.Ttl < 1 {
		log.Println("bad ttl:", r.Ttl)
		return ErrBadTTL
	}
	return nil
}

func GoodName(name string) bool {
	if len(name) < 1 || len(name) > 255 {
		log.Println("bad name:", len(name), name)
		return false
	}
	name = strings.ToLower(name)
	validMap := map[rune]bool{
		'a': true,
		'b': true,
		'c': true,
		'd': true,
		'e': true,
		'f': true,
		'g': true,
		'h': true,
		'i': true,
		'j': true,
		'k': true,
		'l': true,
		'm': true,
		'n': true,
		'o': true,
		'p': true,
		'q': true,
		'r': true,
		's': true,
		't': true,
		'u': true,
		'v': true,
		'w': true,
		'x': true,
		'y': true,
		'z': true,
		'-': true,
		'.': true,
	}
	for _, n := range name {
		_, ok := validMap[n]
		if !ok {
			return false
		}
	}
	return true
}

func AppendDot(name string) string {
	if name[len(name)-1] != '.' {
		name = name + "."
	}
	return name
}

func TypeStrToInt(tp string) uint16 {
	tp = strings.ToUpper(tp)
	switch tp {
	case "A":
		return dns.TypeA
	case "AAAA":
		return dns.TypeAAAA
	case "NS":
		return dns.TypeNS
	case "CNAME":
		return dns.TypeCNAME
	case "TXT":
		return dns.TypeTXT
	case "CAA":
		return dns.TypeCAA
	case "SVCB":
		return dns.TypeSVCB
	case "SOA":
		return dns.TypeSOA
	}
	return dns.TypeNone
}
