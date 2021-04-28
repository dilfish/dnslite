// Copyright 2018 Sean.ZH

package main

import (
	"log"
	"net/http"

	"github.com/dilfish/dnslite"
	"github.com/miekg/dns"
)

// UpDNS create new dns service
func UpDNS() {
	mux := dnslite.CreateDNSMux()
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", mux.ServeDNS)
	err := server.ListenAndServe()
	panic(err)
}

// UpDoT
func UpDoT() {
	cert := "/etc/letsencrypt/live/dilfish.dev-0001/fullchain.pem"
	key := "/etc/letsencrypt/live/dilfish.dev-0001/privkey.pem"
	var h dnslite.Handler
	err := dns.ListenAndServeTLS(":853", cert, key, &h)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go UpDoT()
	go UpDNS()
	mux := dnslite.CreateHTTPMux()
	err := http.ListenAndServe(":8085", mux)
	panic(err)
}
