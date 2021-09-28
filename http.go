// Package dnslite HTTP Api is for debug or human reading
// not compatible with RFC 8484

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/miekg/dns"
)

type HTTPHandler struct {
	M DataManagerI
}

type DNSResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  []DNSRecord `json:"result"`
}

func TypeStrToInt(tp string) uint16 {
	switch tp {
	case "A":
		return dns.TypeA
	case "AAAA":
		return dns.TypeAAAA
	}
	return dns.TypeNone
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	domain := r.Form["domain"]
	tp := r.Form["type"]
	if len(domain) < 1 || len(tp) < 1 {
		w.Write([]byte("please input domain and type"))
		return
	}
	t := TypeStrToInt(tp[0])
	name := domain[0]
	if name[len(name)-1] != '.' {
		name = name + "."
	}
	log.Println("request is", name, t)
	rs, err := h.M.Find(name, t)
	if err != nil {
		log.Println("find error:", name, t, err)
		msg := new(dns.Msg)
		msg.Id = dns.Id()
		msg.RecursionDesired = true
		msg.Question = make([]dns.Question, 1)
		msg.Question[0] = dns.Question{Name: name, Qtype: t, Qclass: dns.ClassINET}
		dret, err := GetDataFromRealDNS(msg)
		if err != nil {
			log.Println("get proxy error:", err)
			w.Write([]byte("get proxy error"))
			return
		}
		w.Write([]byte(dret.String()))
		return
	}
	var ret DNSResult
	ret.Result = rs
	bt, _ := json.Marshal(ret)
	w.Write(bt)
}

func NewHTTPHandler(conf *Config) {
	var srv http.Server
	var h HTTPHandler
	h.M = NewMongoClient(&conf.MongoClientConfig)
	if conf.UsingMemDB {
		h.M = NewMemDB()
	}
	addr := ":" + strconv.FormatInt(int64(conf.Port), 10)
	log.Println("http dns serves on:", addr)
	srv.Addr = addr
	srv.Handler = &h
	err := srv.ListenAndServe()
	if err != nil {
		log.Println("http serve error:", addr, err)
	}
}
