package main

import (
	"context"
	"log"

	"github.com/miekg/dns"
)

func GetDataFromRealDNS(req *dns.Msg, withTLS bool) (int, *dns.Msg, error) {
	dnsMap := map[string]bool{
		"119.29.29.29":    false, // tencent
		"114.114.114.114": false, // 114
		// "1.1.1.1":         true,  // cloudflare
		"8.8.8.8": false, // google
		"dot.pub": true,  // dot.pub
	}
	for addr, tls := range dnsMap {
		if withTLS && !tls {
			continue
		}
		c := new(dns.Client)
		a := addr + ":53"
		if withTLS {
			c = &dns.Client{Net: "tcp-tls"}
			a = addr + ":853"
		}
		log.Println("proxy req to", req.Question[0].Name, req.Question[0].Qtype, a)
		ret, _, err := c.ExchangeContext(context.Background(), req, a)
		if err == nil {
			return 0, ret, nil
		} else {
			log.Println("exchange error of:", a, err)
		}
	}
	return dns.RcodeServerFailure, nil, ErrRCode
}
