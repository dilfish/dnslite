// Copyright 2018 Sean.ZH

package dnslite

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

// ErrBadQCount is more than 1 question count
var ErrBadQCount = errors.New("bad question count")

// ErrNotA only supports A record
var ErrNotA = errors.New("type A support only")

// ExtraInfo fills dns cookie and subnet
type ExtraInfo struct {
	Cookie string
	Subnet string
}

func getDNSInfo(r *dns.Msg) (name string, tp uint16, ex ExtraInfo, err error) {
	if len(r.Question) != 1 {
		err = ErrBadQCount
		log.Println("r.question is not 1", len(r.Question))
		return
	}
	log.Println("extra is", r.Extra)
	name = r.Question[0].Name
	tp = r.Question[0].Qtype
	if tp != dns.TypeA && tp != dns.TypeNS && tp != dns.TypeAAAA && tp != dns.TypeTXT && tp != dns.TypeCAA {
		err = ErrNotA
		log.Println("r.q.type is not A", tp)
		return
	}
	if len(r.Extra) > 2 {
		log.Println("extra len", len(r.Extra))
		return
	}
	return
}

func retNS(w dns.ResponseWriter, r *dns.Msg, name string) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	ns := new(dns.NS)
	ns.Hdr.Name = name
	ns.Hdr.Rrtype = dns.TypeNS
	ns.Hdr.Class = dns.ClassINET
	ns.Hdr.Ttl = 60
	ns.Ns = "ns1.dilfish.dev."
	m.Answer = append(m.Answer, ns)
	ns.Ns = "ns2.dilfish.dev."
	m.Answer = append(m.Answer, ns)
	w.WriteMsg(m)
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("we get request from", w.RemoteAddr(), r.Question)
	log.Println("flags are, auth:", r.Authoritative, "truncated:", r.Truncated, "recursiondesired:", r.RecursionDesired, "recursionavaliable:", r.RecursionAvailable, "ad:", r.AuthenticatedData, "cd:", r.CheckingDisabled)
	m := new(dns.Msg)
	name, tp, _, err := getDNSInfo(r)
	if err != nil {
		log.Println("bad dns info", err)
		return
	}
	m.SetReply(r)
	m.Authoritative = true
	if tp == dns.TypeNS {
		retNS(w, r, name)
		return
	}
	rr, err := GetRecord(name, tp)
	if err != nil {
		log.Println("get record error", name, tp, err)
		return
	}
	if tp == dns.TypeA {
		for _, r := range rr {
			a := new(dns.A)
			a.Hdr.Name = name
			a.Hdr.Rrtype = tp
			a.Hdr.Class = dns.ClassINET
			a.Hdr.Ttl = r.TTL
			a.A = net.ParseIP(r.Value).To4()
			m.Answer = append(m.Answer, a)
		}
	}
	if tp == dns.TypeAAAA {
		for _, r := range rr {
			aaaa := new(dns.AAAA)
			aaaa.Hdr.Name = name
			aaaa.Hdr.Rrtype = tp
			aaaa.Hdr.Class = dns.ClassINET
			aaaa.Hdr.Ttl = r.TTL
			aaaa.AAAA = net.ParseIP(r.Value)
			m.Answer = append(m.Answer, aaaa)
		}
	}
	if tp == dns.TypeTXT {
		for _, r := range rr {
			txt := new(dns.TXT)
			txt.Hdr.Name = name
			txt.Hdr.Rrtype = tp
			txt.Hdr.Class = dns.ClassINET
			txt.Hdr.Ttl = r.TTL
			txt.Txt = strings.Split(r.Value, "\"")
			m.Answer = append(m.Answer, txt)
		}
	}
	if tp == dns.TypeCAA {
		for _, r := range rr {
			caa := new(dns.CAA)
			caa.Hdr.Name = name
			caa.Hdr.Rrtype = tp
			caa.Hdr.Class = dns.ClassINET
			caa.Hdr.Ttl = r.TTL
			caa.Flag = 0
			caa.Tag = "issue"
			caa.Value = r.Value
			m.Answer = append(m.Answer, caa)
		}
	}
	w.WriteMsg(m)
}

// CreateDNSMux create mux for dns like http
func CreateDNSMux() *dns.ServeMux {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", handleRequest)
	return mux
}
