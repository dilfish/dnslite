package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"net/http"
	"time"
)

type TypeRecord struct {
	Type  uint16
	Value string
	Ttl   int
}

// key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	rr := &dns.A{
		Hdr: dns.RR_Header{
			Name:   "baidu.com.",
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    1,
		},
		A: net.ParseIP("1.1.1.1").To4(),
	}
	m.Answer = []dns.RR{rr}
	w.WriteMsg(m)
}

func Handle() error {
	server := &dns.Server{Addr: ":53", Net: "udp4"}
	dns.HandleFunc(".", handleRequest)
	return server.ListenAndServe()
}

func RunHTTP() {
	for {
		err := http.ListenAndServe("127.0.0.1:8083", nil)
		if err != nil {
			time.Sleep(time.Second * 5)
			fmt.Println("listen error", err)
		}
	}
}

func HandleHTTP() {
	http.HandleFunc("/api/add.record", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	http.HandleFunc("/api/del.record", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	http.HandleFunc("/api/list.record", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	go RunHTTP()
}

func main() {
	HandleHTTP()
	err := Handle()
	if err != nil {
		fmt.Println("error is", err)
		panic(err)
	}
}
