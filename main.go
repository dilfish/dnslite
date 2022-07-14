// Copyright 2018 Sean.ZH

package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/miekg/dns"
)

var FlagUsingMemdb = flag.Bool("m", false, "using memory db")
var FlagUsingDnsOverTls = flag.Bool("t", false, "using dns over tls")
var FlagHTTPPort = flag.Int("p", 10083, "http port")
var FlagDebugMode = flag.Bool("d", false, "debug mode")
var FlagAllProxy = flag.Bool("a", false, "all proxy")
var FlagNoneProxy = flag.Bool("n", false, "none proxy")

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
	cert := "/etc/letsencrypt/live/dilfish.dev-0001/fullchain.pem"
	key := "/etc/letsencrypt/live/dilfish.dev-0001/privkey.pem"
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
	log.Println("using mem db:", conf.UsingMemDB)
	if *FlagUsingDnsOverTls {
		go UpDoT(&conf)
	}
	go UpDNS(&conf)
	api := NewApiHandler(&conf)
	if api == nil {
		log.Println("bad api")
		return
	}
	log.Println("listen on:", *FlagHTTPPort)
	err := http.ListenAndServe(":"+strconv.FormatInt(int64(*FlagHTTPPort), 10), api)
	panic(err)
}
