package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var mapLock sync.Mutex

func RunHTTP() {
	for {
		err := http.ListenAndServe("127.0.0.1:8083", nil)
		if err != nil {
			time.Sleep(time.Second * 5)
			fmt.Println("listen error", err)
		}
	}
}

func AddRecord(a AddRecordArgs) error {
	mapLock.Lock()
	defer mapLock.Unlock()
	return nil
}

func DelRecord(d DelRecordArgs) error {
	mapLock.Lock()
	defer mapLock.Unlock()
	return nil
}

func HandleParams(w http.ResponseWriter, r *http.Request, v interface{}) error {
	defer r.Body.Close()
	bt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(""))
		return err
	}
	return json.Unmarshal(bt, v)
	if err != nil {
		w.Write([]byte(""))
		return err
	}
	return nil
}

type AddRecordArgs struct {
	Name  string
	Type  uint16
	Src   string
	Ttl   uint32
	Value string
}

type DelRecordArgs struct {
	Name string
	Type uint16
	Src  string
}

func HandleHTTP() {
	http.HandleFunc("/api/add.record", func(w http.ResponseWriter, r *http.Request) {
		var a AddRecordArgs
		err := HandleParams(w, r, &a)
		if err != nil {
			return
		}
		ret := AddRecord(a)
		bt, _ := json.Marshal(ret)
		w.Write(bt)
		return
	})
	http.HandleFunc("/api/del.record", func(w http.ResponseWriter, r *http.Request) {
		var d DelRecordArgs
		err := HandleParams(w, r, &d)
		if err != nil {
			return
		}
		ret := DelRecord(d)
		bt, _ := json.Marshal(ret)
		w.Write(bt)
		return
	})
	go RunHTTP()
}
