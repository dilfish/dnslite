// Copyright 2021 Sean.ZH

package dnslite

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *ApiHandler) DelRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("del record bad method:", r.Method)
		w.Write(a.BadMethodMsg)
		return
	}
	var record DNSRecord
	err := a.UnjsonRequest(r, &record)
	if err != nil {
		log.Println("unjson req error:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	if record.Name == "" || record.Type == 0 {
		log.Println("empty id for del record")
		w.Write(a.BadRequestMsg)
		return
	}
	a.DB.Del(record.Name, record.Type)
	bt, _ := json.Marshal(record)
	w.Write(bt)
}
