// Copyright 2018 Sean.ZH

package main

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
	DB               DataManagerI
	OkMsg            []byte
	BadRequestMsg    []byte
	NotSupportedType []byte
	BadRecordValue   []byte
	DBErrMsg         []byte
	BadMethodMsg     []byte
	TypeConflictMsg  []byte
}

func NewApiHandler(conf *Config) *ApiHandler {
	var a ApiHandler
	a.DB = GetGlobalDB(conf)
	if a.DB == nil {
		log.Println("db is nil:", conf.UsingMemDB)
		return nil
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/add.record", a.AddRecord)
	mux.HandleFunc("/api/list.record", a.ListRecord)
	mux.HandleFunc("/api/del.record", a.DelRecord)
	mux.HandleFunc("/api/add.rcode", a.AddRcode)
	mux.HandleFunc("/api/del.rcode", a.DelRcode)
	a.Mux = mux
	var ret DNSRecord
	ret.Code = 0
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
	ret.Msg = "db error"
	ieMsg, _ := json.Marshal(ret)
	a.DBErrMsg = ieMsg
	ret.Code = 5
	ret.Msg = "bad method"
	bmMsg, _ := json.Marshal(ret)
	a.BadMethodMsg = bmMsg
	ret.Code = 6
	ret.Msg = "type conflict"
	tcMsg, _ := json.Marshal(ret)
	a.TypeConflictMsg = tcMsg
	return &a
}

func (a *ApiHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	a.Mux.ServeHTTP(rw, req)
}
