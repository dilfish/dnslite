package dnslite

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var errBadQCount = errors.New("bad question count")
var errNotA = errors.New("a support only")

func getIP(addr string) (string, error) {
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

func getDNSInfo(r *dns.Msg) (name string, tp uint16, err error) {
	if len(r.Question) != 1 {
		err = errBadQCount
		return
	}
	name = r.Question[0].Name
	tp = r.Question[0].Qtype
	if tp != dns.TypeA {
		err = errNotA
		return
	}
	return
}

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println(time.Now(), w.RemoteAddr(), r.Question)
	m := new(dns.Msg)
	name, tp, err := getDNSInfo(r)
	if err != nil {
		fmt.Println(time.Now(), "bad dns info", r, err)
		return
	}
	src, err := getIP(w.RemoteAddr().String())
	if err != nil {
		fmt.Println("bad remoteaddr", w.RemoteAddr())
		return
	}
	m.SetReply(r)
	m.Authoritative = true
	rr, err := getRecord(name, tp)
	if err != nil {
		fmt.Println("get record", name, src, tp, err)
		return
	}
	for _, r := range rr {
		a := new(dns.A)
		a.Hdr.Name = name
		a.Hdr.Rrtype = tp
		a.Hdr.Class = dns.ClassINET
		a.Hdr.Ttl = r.TTL
		a.A = net.ParseIP(r.Value).To4()
		m.Answer = append(m.Answer, a)
	}
	w.WriteMsg(m)
}

func Handle() error {
	server := &dns.Server{Addr: ":53", Net: "udp4"}
	dns.HandleFunc(".", handleRequest)
	return server.ListenAndServe()
}
