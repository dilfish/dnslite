// Copyright 2021 Sean.ZH

package dnslite

import (
	"encoding/json"
	"log"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *ApiHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
	cf, ok := TypeHandlerList[record.Type]
	if !ok {
		w.Write(a.NotSupportedType)
		return
	}
	err = CommonCheck(&record)
	if err != nil {
		w.Write(a.BadRecordValue)
		return
	}
	err = cf.CheckRecord(&record)
	if err != nil {
		w.Write(a.BadRecordValue)
		return
	}
	record.Id = primitive.NewObjectID()
	err = a.DB.Insert(record)
	if err != nil {
		w.Write(a.InsertErrMsg)
		return
	}
	record.Msg = "ok"
	bt, _ := json.Marshal(record)
	w.Write(bt)
}
