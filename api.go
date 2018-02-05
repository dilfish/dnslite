package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var ErrValExists = errors.New("no default line")
var ErrNoSuchVal = errors.New("no such value")

type TypeRecord struct {
	Value string
	Ttl   uint32
}

// key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

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

func GetRecord(name, src string, tp uint16) ([]TypeRecord, error) {
	mapLock.Lock()
	defer mapLock.Unlock()
	defKey := name + strconv.Itoa(int(tp))
	realKey := name + strconv.Itoa(int(tp)) + src
	vs, ok := RecordMap[realKey]
	if ok == true {
		return vs, nil
	}
	vs, ok = RecordMap[defKey]
	if ok == true {
		return vs, nil
	}
	return nil, ErrNoSuchVal
}

func AddRecord(a RecordArgs) error {
	mapLock.Lock()
	key := a.Name + strconv.Itoa(int(a.Type)) + a.Src
	var val TypeRecord
	val.Ttl = a.Ttl
	val.Value = a.Value
	vs, ok := RecordMap[key]
	if ok == false {
		RecordMap[key] = []TypeRecord{val}
	} else {
		for _, v := range vs {
			if v.Value == a.Value {
				return ErrValExists
			}
		}
		vs = append(vs, val)
		RecordMap[key] = vs
	}
	defer mapLock.Unlock()
	return nil
}

func DelRecord(d RecordArgs) error {
	mapLock.Lock()
	key := d.Name + strconv.Itoa(int(d.Type)) + d.Src
	vs, ok := RecordMap[key]
	if ok == false {
		return ErrNoSuchVal
	}
	nv := make([]TypeRecord, 0)
	for _, val := range vs {
		if val.Value != d.Value {
			nv = append(nv, val)
		}
	}
	RecordMap[key] = nv
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

type RecordArgs struct {
	Name  string
	Type  uint16
	Src   string
	Ttl   uint32
	Value string
}

func HandleHTTP() {
	http.HandleFunc("/api/add.record", func(w http.ResponseWriter, r *http.Request) {
		var a RecordArgs
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
		var d RecordArgs
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
