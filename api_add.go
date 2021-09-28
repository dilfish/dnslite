// Copyright 2021 Sean.ZH

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/miekg/dns"
)

func (a *ApiHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Write(a.BadMethodMsg)
		log.Println("bad method for add record:", r.Method)
		return
	}
	var record DNSRecord
	err := a.UnjsonRequest(r, &record)
	if err != nil {
		log.Println("unjson req error:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	if record.Name[len(record.Name)-1] != '.' {
		record.Name = record.Name + "."
	}
	cf, ok := TypeHandlerList[record.Type]
	if !ok {
		log.Println("not supported type:", record.Type)
		w.Write(a.NotSupportedType)
		return
	}
	err = CommonCheck(&record)
	if err != nil {
		log.Println("failed of common check:", err)
		w.Write(a.BadRecordValue)
		return
	}
	err = cf.CheckRecord(&record)
	if err != nil {
		w.Write(a.BadRecordValue)
		log.Println("check record:", err)
		return
	}
	conflictList := []uint16{dns.TypeNS, dns.TypeCNAME}
	for _, tp := range conflictList {
		// more than one ns records does not conflicts
		if tp == dns.TypeNS && record.Type == dns.TypeNS {
			continue
		}
		ret, err := a.DB.Find(record.Name, tp)
		if err == nil && len(ret) > 0 {
			log.Println("conflict", dns.Type(record.Type), "with", dns.Type(tp))
			w.Write(a.TypeConflictMsg)
			return
		}
	}
	err = a.DB.Insert(record)
	if err != nil {
		log.Println("db insert error:", err)
		w.Write(a.DBErrMsg)
		return
	}
	record.Msg = "ok"
	bt, _ := json.Marshal(record)
	w.Write(bt)
}
