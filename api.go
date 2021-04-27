// Copyright 2018 Sean.ZH

package dnslite

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/miekg/dns"
)

var errValExists = errors.New("no default line")
var errNoSuchVal = errors.New("no such value")
var errBadName = errors.New("bad name")
var errBadType = errors.New("bad type")
var errBadTTL = errors.New("bad ttl")
var errBadValue = errors.New("bad value")

// TypeRecord reprents a record value and ttl
type TypeRecord struct {
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

// RecordInfo fills typerecord and name of a domain
type RecordInfo struct {
	TypeRecord
	Name string `json:"name"`
}

// RecordMap holds key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

var mapLock sync.Mutex

func listRecord() []RecordInfo {
	var rs []RecordInfo
	var r RecordInfo
	mapLock.Lock()
	defer mapLock.Unlock()
	for k, vs := range RecordMap {
		r.Name = k
		for _, v := range vs {
			r.Value = v.Value
			r.TTL = v.TTL
			rs = append(rs, r)
		}
	}
	return rs
}

// GetRecord returns an certain name and type value to client
func GetRecord(name string, tp uint16) ([]TypeRecord, error) {
	mapLock.Lock()
	defer mapLock.Unlock()
	key := name + strconv.Itoa(int(tp))
	vs, ok := RecordMap[key]
	if ok == true {
		return vs, nil
	}
	return nil, errNoSuchVal
}

func delRecord(d RecordArgs) error {
	if d.Name[len(d.Name)-1] != '.' {
		d.Name = d.Name + "."
	}
	key := d.Name + strconv.Itoa(int(d.Type))
	mapLock.Lock()
	defer mapLock.Unlock()
	v, ok := RecordMap[key]
	if ok == false {
		return nil
	}
	list := make([]TypeRecord, 0)
	for _, item := range v {
		if item.Value != d.Value {
			list = append(list, item)
		}
	}
	RecordMap[key] = list
	return nil
}

func addRecord(a RecordArgs) error {
	if a.Name == "" {
		return errBadName
	}
	if a.Name[len(a.Name)-1] != '.' {
		a.Name = a.Name + "."
	}
	if !isSupportedType(a.Type) {
		log.Println("we does not support type:", a.Type)
		return errBadType
	}
	if a.TTL > 600 || a.TTL < 1 {
		return errBadTTL
	}
	switch a.Type {
		case dns.TypeA:
			ip := net.ParseIP(a.Value)
			if ip == nil {
				return errBadValue
			}
			ip = ip.To4()
			if ip == nil {
				return errBadValue
			}
		case dns.TypeAAAA:
			ip := net.ParseIP(a.Value)
			if ip == nil {
				return errBadValue
			}
			ip = ip.To4()
			if ip != nil {
				return errBadValue
			}
		// no need to check cname, ns, caa and txt
	}
	key := a.Name + strconv.Itoa(int(a.Type))
	var val TypeRecord
	val.TTL = a.TTL
	val.Value = a.Value
	mapLock.Lock()
	defer mapLock.Unlock()
	v, ok := RecordMap[key]
	if ok == false {
		RecordMap[key] = []TypeRecord{val}
	} else {
		for _, item := range v {
			if item.Value == a.Value {
				return errBadValue
			}
		}
		v = append(v, val)
		RecordMap[key] = v
	}
	return nil
}

func handleParams(w http.ResponseWriter, r *http.Request, v interface{}) error {
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

// RecordArgs would be send when calls api
type RecordArgs struct {
	Name  string `json:"name"`
	Type  uint16 `json:"type"`
	TTL   uint32 `json:"ttl"`
	Value string `json:"value"`
}

// RecordRet is the response
type RecordRet struct {
	Err int    `json:"err"`
	Msg string `json:"msg"`
}

// del record ignore ttl argument
func httpDelRecord(w http.ResponseWriter, r *http.Request) {
	var ret RecordRet
	var d RecordArgs
	ret.Msg = "ok"
	err := handleParams(w, r, &d)
	if err != nil {
		return
	}
	err = delRecord(d)
	if err != nil {
		ret.Err = 2
		ret.Msg = err.Error()
	}
	bt, _ := json.Marshal(ret)
	w.Write(bt)
	return
}

func httpAddRecord(w http.ResponseWriter, r *http.Request) {
	var ret RecordRet
	var a RecordArgs
	ret.Msg = "ok"
	err := handleParams(w, r, &a)
	if err != nil {
		return
	}
	err = addRecord(a)
	if err != nil {
		ret.Err = 1
		ret.Msg = err.Error()
	}
	bt, _ := json.Marshal(ret)
	w.Write(bt)
}

func httpListRecord(w http.ResponseWriter, r *http.Request) {
	l := listRecord()
	bt, _ := json.Marshal(l)
	w.Write(bt)
}

// CreateHTTPMux create http service handler
// for dnslite
func CreateHTTPMux() http.Handler {
	if RecordMap == nil {
		RecordMap = make(map[string][]TypeRecord)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/add.record", httpAddRecord)
	mux.HandleFunc("/api/list.record", httpListRecord)
	mux.HandleFunc("/api/del.record", httpDelRecord)
	return mux
}
