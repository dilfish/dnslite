package dnslite

import (
	"log"
	"strings"
	"github.com/miekg/dns"
)

type TypeHandler interface {
	FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg
	CheckRecord(record *DNSRecord) error
}

var TypeHandlerList = map[uint16]TypeHandler{
	dns.TypeA: &AHandler{},
	dns.TypeAAAA: &AAAAHandler{},
	dns.TypeNS: &NSHandler{},
	dns.TypeCNAME: &CNAMEHandler{},
	dns.TypeTXT: &TXTHandler{},
	dns.TypeCAA: &CAAHandler{},
}

func CommonCheck(r *DNSRecord) error {
	if r.Name == "" {
		return ErrBadName
	}
	if r.Name[len(r.Name)-1] != '.' {
		r.Name = r.Name + "."
	}
	if r.Ttl > 600 || r.Ttl < 1 {
		return ErrBadTTL
	}
	return nil
}

func GoodName(name string) bool {
	if len(name) < 1 || len(name) > 255 {
		log.Println("bad name", len(name), name)
		return false
	}
	name = strings.ToLower(name)
	validMap := map[rune]bool {
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
