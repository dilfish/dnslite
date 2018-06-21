package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var ErrValExists = errors.New("no default line")
var ErrNoSuchVal = errors.New("no such value")
var ErrBadName = errors.New("bad name")
var ErrBadType = errors.New("bad type")
var ErrBadTtl = errors.New("bad ttl")
var ErrBadValue = errors.New("bad value")

type TypeRecord struct {
	Value string `json:"value"`
	Ttl   uint32 `json:"ttl"`
}

type RecordInfo struct {
	TypeRecord
	Name string `json:"name"`
}

// key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

var mapLock sync.Mutex

func RunHTTP() {
	for {
		err := http.ListenAndServe("127.0.0.1:8083", nil)
		if err != nil {
			time.Sleep(time.Second * 5)
			log.Println("listen error", err)
		}
	}
}

func ListRecord() []RecordInfo {
	rs := make([]RecordInfo, 0)
	var r RecordInfo
	mapLock.Lock()
	defer mapLock.Unlock()
	for k, vs := range RecordMap {
		r.Name = k
		for _, v := range vs {
			r.Value = v.Value
			r.Ttl = v.Ttl
			rs = append(rs, r)
		}
	}
	return rs
}

func GetRecord(name string, tp uint16) ([]TypeRecord, error) {
	mapLock.Lock()
	defer mapLock.Unlock()
	key := name + strconv.Itoa(int(tp))
	vs, ok := RecordMap[key]
	if ok == true {
		return vs, nil
	}
	return nil, ErrNoSuchVal
}

func AddRecord(a RecordArgs) error {
	if a.Name == "" {
		return ErrBadName
	}
	if a.Name[len(a.Name)-1] != "." {
		a.Name = a.Name + "."
	}
	if a.Type != 1 {
		return ErrBadType
	}
	if a.Ttl > 600 || a.Ttl < 1 {
		return ErrBadTtl
	}
	if net.ParseIP(a.Value) == nil {
		return ErrBadValue
	}
	key := a.Name + strconv.Itoa(int(a.Type))
	var val TypeRecord
	val.Ttl = a.Ttl
	val.Value = a.Value
	mapLock.Lock()
	defer mapLock.Unlock()
	RecordMap[key] = []TypeRecord{val}
	return nil
}

func HandleParams(w http.ResponseWriter, r *http.Request, v interface{}) error {
	defer r.Body.Close()
	bt, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("err is:" + err.Error()))
		return err
	}
	err = json.Unmarshal(bt, v)
	if err != nil {
		w.Write([]byte("err is:" + err.Error()))
		return err
	}
	return nil
}

type RecordArgs struct {
	Name  string `json:"name"`
	Type  uint16 `json:"type"`
	Ttl   uint32 `json:"ttl"`
	Value string `json:"value"`
}

type RecordRet struct {
	Err int    `json:"err"`
	Msg string `json:"msg"`
}

func HandleHTTP() {
	http.HandleFunc("/api/add.record", func(w http.ResponseWriter, r *http.Request) {
		var ret RecordRet
		var a RecordArgs
		ret.Msg = "ok"
		err := HandleParams(w, r, &a)
		if err != nil {
			return
		}
		err = AddRecord(a)
		if err != nil {
			ret.Err = 1
			ret.Msg = err.Error()
		}
		bt, _ := json.Marshal(ret)
		w.Write(bt)
		return
	})
	http.HandleFunc("/api/list.record", func(w http.ResponseWriter, r *http.Request) {
		l := ListRecord()
		bt, _ := json.Marshal(l)
		w.Write(bt)
		return
	})
	go RunHTTP()
}
