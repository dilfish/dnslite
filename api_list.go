// Copyright 2021 Sean.ZH

package dnslite

import (
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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
	var ret []DNSRecord
	err = a.DB.Find(bson.M{"name": record.Name, "type": record.Type}, &ret)
	if err != nil {
		log.Println("insert error:", err)
		w.Write(a.DBErrMsg)
		return
	}
	bt, _ := json.Marshal(ret)
	w.Write(bt)
}
