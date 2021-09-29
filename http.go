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

func MsgToRecord(msg *dns.Msg) []DNSRecord {
	list := make([]DNSRecord, 0)
	for _, a := range msg.Answer {
		t := a.Header().Rrtype
		_, ok := TypeHandlerList[t]
		if !ok {
			log.Println("no such type to parse:", t)
			continue
		}
		record := TypeHandlerList[t].RRToRecord(a)
		record.Name = a.Header().Name
		record.Type = a.Header().Rrtype
		list = append(list, record)
	}
	for _, a := range msg.Ns {
		t := a.Header().Rrtype
		_, ok := TypeHandlerList[t]
		if !ok {
			log.Println("no such type to parse:", t)
			continue
		}
		record := TypeHandlerList[t].RRToRecord(a)
		record.Name = a.Header().Name
		record.Type = a.Header().Rrtype
		list = append(list, record)
	}
	for _, a := range msg.Extra {
		t := a.Header().Rrtype
		_, ok := TypeHandlerList[t]
		if !ok {
			log.Println("no such type to parse:", t)
			continue
		}
		record := TypeHandlerList[t].RRToRecord(a)
		record.Name = a.Header().Name
		record.Type = a.Header().Rrtype
		list = append(list, record)
	}
	return list
}

func HTTPProxy(name string, tp uint16) ([]DNSRecord, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{Name: name, Qtype: tp, Qclass: dns.ClassINET}
	dret, err := GetDataFromRealDNS(msg)
	if err != nil {
		return nil, err
	}
	if *FlagDebugMode {
		log.Println("http get", dret)
	}
	return MsgToRecord(dret), nil
}

func (h *HTTPHandler) getRecord(domain, tp []string) DNSResult {
	var result DNSResult
	if len(domain) < 1 || len(tp) < 1 {
		log.Println("bad domain or type:", domain, tp)
		result.Code = 1
		result.Message = "please input domain and type"
		return result
	}
	t := TypeStrToInt(tp[0])
	if t == dns.TypeNone {
		log.Println("bad request type:", tp[0])
		result.Code = 2
		result.Message = "request type not supported:" + tp[0]
		return result
	}
	name := domain[0]
	if name[len(name)-1] != '.' {
		name = name + "."
	}
	log.Println("request is", name, t)
	rs, err := h.M.Find(name, t)
	if err != nil || len(rs) == 0 {
		log.Println("find error:", name, t, err)
		rs, err = HTTPProxy(name, t)
		if err != nil {
			log.Println("get proxy error:", err)
			result.Code = 3
			result.Message = "get proxy error:" + err.Error()
			return result
		}
	}
	result.Result = rs
	return result
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	domain := r.Form["domain"]
	tp := r.Form["type"]
	ret := h.getRecord(domain, tp)
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
