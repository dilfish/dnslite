// Copyright 2018 Sean.ZH

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/dilfish/dnslite"
	"github.com/miekg/dns"
)

var FlagUsingMemdb = flag.Bool("m", false, "using memory db")

// UpDNS create new dns service
func UpDNS(conf *dnslite.Config) {
	h := dnslite.NewHandler(conf)
	err := dns.ListenAndServe(":53", "udp", h)
	if err != nil {
		panic(err)
	}
}

// UpDoT
func UpDoT(conf *dnslite.Config) {
	cert := "./fullchain4.pem"
	key := "./privkey4.pem"
	h := dnslite.NewHandler(conf)
	err := dns.ListenAndServeTLS(":853", cert, key, h)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	var conf dnslite.Config
	conf.Addr = "mongodb://localhost:27017"
	conf.DB = "dnslite"
	conf.Coll = "records"
	conf.UsingMemDB = *FlagUsingMemdb
	go UpDoT(&conf)
	go UpDNS(&conf)
	api := dnslite.NewApiHandler(&conf)
	err := http.ListenAndServe(":8085", api)
	panic(err)
}
