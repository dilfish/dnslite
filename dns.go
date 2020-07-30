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
var ErrNotA = errors.New("type a support only")

// ExtraInfo fills dns cookie and subnet
type ExtraInfo struct {
	Cookie string
	Subnet string
}

func isSupportedType(tp uint16) bool {
	switch tp {
	case dns.TypeA: fallthrough
	case dns.TypeAAAA: fallthrough
	case dns.TypeNS: fallthrough
	case dns.TypeTXT: fallthrough
	case dns.TypeCAA: fallthrough
	case dns.TypeCNAME:
		return true
	}
	return false
}

func getDNSInfo(r *dns.Msg) (name string, tp uint16, ex ExtraInfo, err error) {
	if len(r.Question) != 1 {
		err = ErrBadQCount
		log.Println("r.question is not 1", len(r.Question))
		return
	}
	for _, e := range r.Extra {
		x, ok := e.(*dns.OPT)
		if ok {
			log.Println("subnet and udp size:", x.Option, x.UDPSize())
		} else {
			log.Println("Unkown extra is:", e)
		}
	}
	name = r.Question[0].Name
	tp = r.Question[0].Qtype
	if !isSupportedType(tp) {
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

func fillHdr(hdr *dns.RR_Header, name string, tp uint16, ttl uint32) {
	hdr.Name = name
	hdr.Ttl = ttl
	hdr.Class = dns.ClassINET
	hdr.Rrtype = tp
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("we get request from", w.RemoteAddr(), r.Question)
	log.Println("flags are, auth:", r.Authoritative, ", trunc:", r.Truncated, ", recur desired:", r.RecursionDesired, ", recur avail:", r.RecursionAvailable, "ad:", r.AuthenticatedData, "cd:", r.CheckingDisabled)
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
	// we write back soa
	if err == errNoSuchVal {
		a := new(dns.SOA)
		fillHdr(&a.Hdr, name, tp, 600)
		a.Serial = 2012143034
		a.Refresh = 300
		a.Expire = 2592000
		a.Retry = 300
		a.Mbox = "i.at.dilfish.dev."
		a.Minttl = 100
		a.Ns = "ns1.dilfish.dev."
		m.Answer = append(m.Answer, a)
		w.WriteMsg(m)
		return
	}
	if err != nil {
		log.Println("get record error", name, tp, err)
		return
	}
	if tp == dns.TypeA {
		for _, r := range rr {
			a := new(dns.A)
			fillHdr(&a.Hdr, name, tp, r.TTL)
			a.A = net.ParseIP(r.Value).To4()
			m.Answer = append(m.Answer, a)
		}
	}
	if tp == dns.TypeAAAA {
		for _, r := range rr {
			aaaa := new(dns.AAAA)
			fillHdr(&aaaa.Hdr, name, tp, r.TTL)
			aaaa.AAAA = net.ParseIP(r.Value)
			m.Answer = append(m.Answer, aaaa)
		}
	}
	if tp == dns.TypeTXT {
		for _, r := range rr {
			txt := new(dns.TXT)
			fillHdr(&txt.Hdr, name, tp, r.TTL)
			txt.Txt = strings.Split(r.Value, "\"")
			m.Answer = append(m.Answer, txt)
		}
	}
	if tp == dns.TypeCAA {
		for _, r := range rr {
			caa := new(dns.CAA)
			fillHdr(&caa.Hdr, name, tp, r.TTL)
			caa.Flag = 0
			caa.Tag = "issue"
			caa.Value = r.Value
			m.Answer = append(m.Answer, caa)
		}
	}
	if tp == dns.TypeCNAME {
		for _, r := range rr {
			cname := new(dns.CNAME)
			fillHdr(&cname.Hdr, name, tp, r.TTL)
			cname.Target = r.Value
			m.Answer = append(m.Answer, cname)
		}
	}
	w.WriteMsg(m)
}

// CreateDNSMux create mux for dns like http
func CreateDNSMux() *dns.ServeMux {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := dns.NewServeMux()
	mux.HandleFunc(".", handleRequest)
	return mux
}
