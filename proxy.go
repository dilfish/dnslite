package main

import (
	"log"

	"github.com/miekg/dns"
)

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
