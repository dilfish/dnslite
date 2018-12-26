// Copyright 2018 Sean.ZH

package main

import (
	"github.com/dilfish/dnslite"
	"github.com/miekg/dns"
	"net/http"
)

// UpDNS create new dns service
func UpDNS() {
	mux := dnslite.CreateDNSMux()
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", mux.ServeDNS)
	err := server.ListenAndServe()
	panic(err)
}

func main() {
	go UpDNS()
	mux := dnslite.CreateHTTPMux()
	err := http.ListenAndServe(":8085", mux)
	panic(err)
}
