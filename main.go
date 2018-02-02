package main

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
)

type TypeRecord struct {
	Type  uint16
	Value string
	Ttl   int
}

// key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

func GetIP(addr string) (string, error) {
	arr := strings.Split(addr, ":")
	if len(arr) != 2 {
		fmt.Println("bad format addr", addr)
		return "", errors.New("bad format addr")
	}
	ip := net.ParseIP(arr[0])
	if ip == nil {
		fmt.Println("ip is", arr[0])
		return "", errors.New("bad format ip")
	}
	if ip = ip.To4(); ip == nil {
		fmt.Println("ip is", arr[0])
		return "", errors.New("bad format ip")
	}
	return ip.String(), nil
}

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
	fmt.Println(w.RemoteAddr())
}

func Handle() error {
	server := &dns.Server{Addr: ":53", Net: "udp4"}
	dns.HandleFunc(".", handleRequest)
	return server.ListenAndServe()
}

func main() {
	HandleHTTP()
	err := Handle()
	if err != nil {
		fmt.Println("error is", err)
		panic(err)
	}
}
