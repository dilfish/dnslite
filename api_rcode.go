package main

import (
	"log"
	"net/http"

	"github.com/miekg/dns"
)

type RcodeArgs struct {
	Name  string `json:"name"`
	Rcode int    `json:"rcode"`
}

func (a *ApiHandler) AddRcode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(a.BadMethodMsg)
		log.Println("bad method for add rcode:", r.Method)
		return
	}
	var rcode RcodeArgs
	err := a.UnjsonRequest(r, &rcode)
	if err != nil {
		log.Println("unjson req error:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	if rcode.Name == "" || rcode.Rcode == 0 {
		log.Println("bad rcode:", rcode.Name, rcode.Rcode)
		w.Write(a.BadRequestMsg)
		return
	}
	rcode.Name = dns.Fqdn(rcode.Name)
	RcodeMap[rcode.Name] = rcode.Rcode
	log.Println("add rcode:", rcode.Name, rcode.Rcode)
	w.Write(a.OkMsg)
}

func (a *ApiHandler) DelRcode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(a.BadMethodMsg)
		log.Println("bad method for add rcode:", r.Method)
		return
	}
	var rcode RcodeArgs
	err := a.UnjsonRequest(r, &rcode)
	if err != nil {
		log.Println("unjson req error:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	if rcode.Name == "" {
		log.Println("bad rcode:", rcode.Name, rcode.Rcode)
		w.Write(a.BadRequestMsg)
		return
	}
	delete(RcodeMap, rcode.Name)
	log.Println("delete rcode:", rcode.Name)
	w.Write(a.OkMsg)
}
