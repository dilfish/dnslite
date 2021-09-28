// Copyright 2021 Sean.ZH

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *ApiHandler) ListRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("list record bad method:", r.Method)
		w.Write(a.BadMethodMsg)
		return
	}
	var record DNSRecord
	err := a.UnjsonRequest(r, &record)
	if err != nil {
		log.Println("bad request:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	// for good check
	record.Ttl = 100
	err = CommonCheck(&record)
	if err != nil {
		log.Println("bad common check:", record)
		w.Write(a.BadRecordValue)
		return
	}
	ret, err := a.DB.Find(record.Name, record.Type)
	if err != nil {
		log.Println("find error:", err)
		w.Write(a.DBErrMsg)
		return
	}
	// empty result, make an empty slice
	if len(ret) == 0 {
		ret = make([]DNSRecord, 0)
	}
	bt, _ := json.Marshal(ret)
	w.Write(bt)
}
