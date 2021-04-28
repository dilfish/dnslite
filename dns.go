// Copyright 2018 Sean.ZH

package dnslite

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

// ExtraInfo fills dns cookie and subnet
type ExtraInfo struct {
	Cookie string
	Subnet string
}

func isSupportedType(tp uint16) bool {
	switch tp {
	case dns.TypeA:
		fallthrough
	case dns.TypeAAAA:
		fallthrough
	case dns.TypeNS:
		fallthrough
	case dns.TypeTXT:
		fallthrough
	case dns.TypeCAA:
		fallthrough
	case dns.TypeCNAME:
		return true
	case dns.TypePTR:
		return true
	}
	return false
}

func getDNSInfo(r *dns.Msg) (name string, tp dns.Type, ex ExtraInfo, err error) {
	if len(r.Question) != 1 {
		err = ErrBadQCount
		log.Println("r.question is not 1", len(r.Question))
		return
	}
	for _, e := range r.Extra {
		x, ok := e.(*dns.OPT)
		if ok {
			log.Println("subnet and udp size:", len(x.Option), x.UDPSize())
			for _, o := range x.Option {
				switch o.Option() {
				case dns.EDNS0SUBNET:
					log.Println("we get subnet:")
					sub, ok := o.(*dns.EDNS0_SUBNET)
					if ok {
						log.Println("sub info:", sub.String())
					}
				case dns.EDNS0COOKIE:
					log.Println("we get cookie")
					cookie, ok := o.(*dns.EDNS0_COOKIE)
					if ok {
						log.Println("cookie info:", cookie.String())
					}
				}
			}
		} else {
			log.Println("Unkown extra is:", e.String())
		}
	}
	name = r.Question[0].Name
	tp = dns.Type(r.Question[0].Qtype)
	if !isSupportedType(uint16(tp)) {
		err = ErrNotSupported
		log.Println("request type is not supported:", dns.Type(r.Question[0].Qtype))
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
	var ns1, ns2 dns.NS
	ns1.Hdr.Name = name
	ns2.Hdr.Name = name
	ns1.Hdr.Rrtype = dns.TypeNS
	ns2.Hdr.Rrtype = dns.TypeNS
	ns1.Hdr.Class = dns.ClassINET
	ns2.Hdr.Class = dns.ClassINET
	ns1.Hdr.Ttl = 60
	ns2.Hdr.Ttl = 60
	ns1.Ns = "ns1.dilfish.dev."
	m.Answer = append(m.Answer, &ns1)
	ns2.Ns = "ns2.dilfish.dev."
	m.Answer = append(m.Answer, &ns2)
	w.WriteMsg(m)
}

func fillHdr(hdr *dns.RR_Header, name string, tp uint16, ttl uint32) {
	hdr.Name = name
	hdr.Ttl = ttl
	hdr.Class = dns.ClassINET
	hdr.Rrtype = tp
}

func GetDataFromRealDNS(req *dns.Msg) (*dns.Msg, error) {
	dnsList := []string{
		"119.29.29.29",    // tencent
		"114.114.114.114", // 114
		"1.1.1.1",         // cloudflare
		"8.8.8.8",         // google
	}
	for _, d := range dnsList {
		c := new(dns.Client)
		ret, _, err := c.Exchange(req, d+":53")
		if err == nil {
			return ret, nil
		} else {
			log.Println("exchange error of:", d, err)
		}
	}
	return nil, ErrNoGoodServers
}

type Handler struct{}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("remote addr:", w.RemoteAddr())
	log.Println("flags are, auth:", r.Authoritative, ", trunc:", r.Truncated, ", recur desired:", r.RecursionDesired, ", recur avail:", r.RecursionAvailable, "ad:", r.AuthenticatedData, "cd:", r.CheckingDisabled)
	m := new(dns.Msg)
	name, tp, _, err := getDNSInfo(r)
	if err != nil {
		log.Println("bad dns info", err)
		return
	}
	tp16 := uint16(tp)
	log.Println("name and type:", name, tp)
	m.SetReply(r)
	m.Authoritative = true
	if uint16(tp) == dns.TypeNS {
		retNS(w, r, name)
		return
	}
	rr, err := GetRecord(name, uint16(tp))
	// when NON set record is requested, we proxy it to 1.1.1.1
	if err == ErrNoSuchVal {
		r, err := GetDataFromRealDNS(r)
		if err != nil {
			log.Println("exchange error:", err)
			return
		}
		w.WriteMsg(r)
		return
	}
	if err != nil {
		log.Println("get record error", name, tp, err)
		return
	}
	log.Println("hit cache", name, tp)
	if tp16 == dns.TypeA {
		for _, r := range rr {
			a := new(dns.A)
			fillHdr(&a.Hdr, name, tp16, r.TTL)
			a.A = net.ParseIP(r.Value).To4()
			m.Answer = append(m.Answer, a)
		}
	}
	if tp16 == dns.TypeAAAA {
		for _, r := range rr {
			aaaa := new(dns.AAAA)
			fillHdr(&aaaa.Hdr, name, tp16, r.TTL)
			aaaa.AAAA = net.ParseIP(r.Value)
			m.Answer = append(m.Answer, aaaa)
		}
	}
	if tp16 == dns.TypeTXT {
		for _, r := range rr {
			txt := new(dns.TXT)
			fillHdr(&txt.Hdr, name, tp16, r.TTL)
			txt.Txt = strings.Split(r.Value, "\"")
			m.Answer = append(m.Answer, txt)
		}
	}
	if tp16 == dns.TypeCAA {
		for _, r := range rr {
			caa := new(dns.CAA)
			fillHdr(&caa.Hdr, name, tp16, r.TTL)
			caa.Flag = 0
			caa.Tag = "issue"
			caa.Value = r.Value
			m.Answer = append(m.Answer, caa)
		}
	}
	if tp16 == dns.TypeCNAME {
		for _, r := range rr {
			cname := new(dns.CNAME)
			fillHdr(&cname.Hdr, name, tp16, r.TTL)
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
	var h Handler
	mux.HandleFunc(".", h.ServeDNS)
	return mux
}
