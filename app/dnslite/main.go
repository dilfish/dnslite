package main

import (
	"github.com/dilfish/dnslite"
	"github.com/miekg/dns"
	"net/http"
)

func UpDNS() {
	mux := dnslite.CreateDNSMux()
	server := &dns.Server{Addr: "53", Net: "udp4"}
	dns.HandleFunc(".", mux.ServeDNS)
	server.ListenAndServe()
}

func main() {
	go UpDNS()
	mux := dnslite.CreateHTTPMux()
	http.ListenAndServe(":8081", mux)
}
