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
var FlagUsingMongo = flag.String("mgo", "", "mongo db addr")
var FlagUsingDnsOverTls = flag.Bool("t", false, "using dns over tls")
var FlagProxyTls = flag.Bool("pt", false, "proxy using tls")
var FlagHTTPPort = flag.Int("p", 10083, "http port")
var FlagDebugMode = flag.Bool("d", false, "debug mode")
var FlagNoneProxy = flag.Bool("n", false, "none proxy")
var FlagRcodeNoRecord = flag.Int("rc", dns.RcodeNameError, "no record error rcode")

// UpDNS create new dns service
func UpDNS(conf *Config, db DataManagerI) {
	h := NewHandler(conf, db)
	err := dns.ListenAndServe(":53", "udp", h)
	if err != nil {
		panic(err)
	}
}

// UpDoT
func UpDoT(conf *Config, db DataManagerI) {
	cert := "/etc/letsencrypt/live/dilfish.dev-0001/fullchain.pem"
	key := "/etc/letsencrypt/live/dilfish.dev-0001/privkey.pem"
	h := NewHandler(conf, db)
	err := dns.ListenAndServeTLS(":853", cert, key, h)
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	var conf Config
	conf.Addr = *FlagUsingMongo
	conf.DB = "dnslite"
	conf.Coll = "records"
	conf.UsingMemDB = *FlagUsingMemdb
	if *FlagRcodeNoRecord == 0 {
		*FlagRcodeNoRecord = dns.RcodeNameError
	}
	log.Println("using mem db:", conf.UsingMemDB)
	db := NewDB(&conf)
	if db == nil {
		log.Println("new db error")
		return
	}
	if *FlagUsingDnsOverTls {
		go UpDoT(&conf, db)
	}
	go UpDNS(&conf, db)
	api := NewApiHandler(&conf, db)
	if api == nil {
		log.Println("bad api")
		return
	}
	log.Println("listen on:", *FlagHTTPPort)
	err := http.ListenAndServe(":"+strconv.FormatInt(int64(*FlagHTTPPort), 10), api)
	panic(err)
}
