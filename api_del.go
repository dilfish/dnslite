// Copyright 2018 Sean.ZH

package dnslite

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func (a *ApiHandler) DelRecord(w http.ResponseWriter, r *http.Request) {
	var record DNSRecord
	err := a.UnjsonRequest(r, &record)
	if err != nil {
		log.Println("unjson req error:", err)
		w.Write(a.BadRequestMsg)
		return
	}
	if record.Id.Hex() == "" {
		w.Write(a.BadRequestMsg)
		return
	}
	a.DB.Del(bson.M{"_id": record.Id})
}
