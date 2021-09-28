// Copyright 2018 Sean.ZH

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/miekg/dns"
)

var FlagUsingMemdb = flag.Bool("m", false, "using memory db")
var FlagUsingDnsOverTls = flag.Bool("t", false, "using dns over tls")
var FlagHTTPPort = flag.Int("p", 0, "http port")

// UpDNS create new dns service
func UpDNS(conf *Config) {
	h := NewHandler(conf)
	err := dns.ListenAndServe(":53", "udp", h)
	if err != nil {
		panic(err)
	}
}

// UpDoT
func UpDoT(conf *Config) {
	cert := "./fullchain4.pem"
	key := "./privkey4.pem"
	h := NewHandler(conf)
	err := dns.ListenAndServeTLS(":853", cert, key, h)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	var conf Config
	conf.Addr = "mongodb://localhost:27017"
	conf.DB = "dnslite"
	conf.Coll = "records"
	conf.UsingMemDB = *FlagUsingMemdb
	if *FlagUsingDnsOverTls {
		go UpDoT(&conf)
	}
	if *FlagHTTPPort != 0 {
		conf.Port = *FlagHTTPPort
		go NewHTTPHandler(&conf)
	}
	go UpDNS(&conf)
	api := NewApiHandler(&conf)
	err := http.ListenAndServe(":8085", api)
	panic(err)
}
