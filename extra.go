package dnslite

import (
	"log"

	"github.com/miekg/dns"
)

func PrintExtra(extra []dns.RR) {
	for _, e := range extra {
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
}
