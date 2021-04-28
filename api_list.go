// Copyright 2018 Sean.ZH

package dnslite

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (a *ApiHandler) ListRecord(w http.ResponseWriter, r *http.Request) {
	var record DNSRecord
	err := a.UnjsonRequest(r, &record)
	if err != nil {
		w.Write(a.BadRequestMsg)
		return
	}
	// for good check
	record.TTL = 100
	err = CommonCheck(&record)
	if err != nil {
		w.Write(a.BadRecordValue)
		return
	}
	var ret []DNSRecord
	err = a.DB.Find(bson.M{"name": record.Name, "type": record.Type}, &ret)
	if err != nil {
		w.Write(a.InsertErrMsg)
		return
	}
	bt, _ := json.Marshal(ret)
	w.Write(bt)
}
