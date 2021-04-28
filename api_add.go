// Copyright 2021 Sean.ZH

package dnslite

import (
	"log"
	"net/http"
)

func (a *ApiHandler) AddRecord(w http.ResponseWriter, r *http.Request) {
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
	err = a.DB.Insert(record)
	if err != nil {
		w.Write(a.InsertErrMsg)
		return
	}
	w.Write(a.OkMsg)
}
