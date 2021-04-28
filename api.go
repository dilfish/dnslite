// Copyright 2018 Sean.ZH

package dnslite

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func (a *ApiHandler) UnjsonRequest(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	bt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read all error:", err)
		return err
	}
	err = json.Unmarshal(bt, v)
	if err != nil {
		log.Println("unjson error:", string(bt), err)
		return err
	}
	return nil
}

type ApiHandler struct {
	Mux              http.Handler
	DB               *MongoClient
	BadRequestMsg    []byte
	NotSupportedType []byte
	BadRecordValue   []byte
	OkMsg            []byte
	InsertErrMsg     []byte
}

func NewApiHandler(conf *MongoClientConfig) *ApiHandler {
	var a ApiHandler
	m := NewMongoClient(conf)
	if m == nil {
		return nil
	}
	a.DB = m
	mux := http.NewServeMux()
	mux.HandleFunc("/api/add.record", a.AddRecord)
	mux.HandleFunc("/api/list.record", a.ListRecord)
	mux.HandleFunc("/api/del.record", a.DelRecord)
	var ret DNSRecord
	ret.Msg = "ok"
	okMsg, _ := json.Marshal(ret)
	a.OkMsg = okMsg
	ret.Code = 1
	ret.Msg = "bad request"
	brMsg, _ := json.Marshal(ret)
	a.BadRequestMsg = brMsg
	ret.Code = 2
	ret.Msg = "not supported type"
	nstMsg, _ := json.Marshal(ret)
	a.NotSupportedType = nstMsg
	ret.Code = 3
	ret.Msg = "bad record value"
	brvMsg, _ := json.Marshal(ret)
	a.BadRecordValue = brvMsg
	ret.Code = 4
	ret.Msg = "insert error"
	ieMsg, _ := json.Marshal(ret)
	a.InsertErrMsg = ieMsg
	return &a
}

func (a *ApiHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.Mux.ServeHTTP(rw, req)
}
