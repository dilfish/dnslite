package dnslite

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

var errValExists = errors.New("no default line")
var errNoSuchVal = errors.New("no such value")
var errBadName = errors.New("bad name")
var errBadType = errors.New("bad type")
var errBadTTL = errors.New("bad ttl")
var errBadValue = errors.New("bad value")

type TypeRecord struct {
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

type recordInfo struct {
	TypeRecord
	Name string `json:"name"`
}

// key: domain + type + fromIP
var RecordMap map[string][]TypeRecord

var mapLock sync.Mutex

func runHTTP() {
	for {
		err := http.ListenAndServe("127.0.0.1:8083", nil)
		if err != nil {
			time.Sleep(time.Second * 5)
			log.Println("listen error", err)
		}
	}
}

func listRecord() []recordInfo {
	rs := make([]recordInfo, 0)
	var r recordInfo
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

func getRecord(name string, tp uint16) ([]TypeRecord, error) {
	mapLock.Lock()
	defer mapLock.Unlock()
	key := name + strconv.Itoa(int(tp))
	vs, ok := RecordMap[key]
	if ok == true {
		return vs, nil
	}
	return nil, errNoSuchVal
}

func addRecord(a recordArgs) error {
	if a.Name == "" {
		return errBadName
	}
	if a.Name[len(a.Name)-1] != '.' {
		a.Name = a.Name + "."
	}
	if a.Type != 1 {
		return errBadType
	}
	if a.TTL > 600 || a.TTL < 1 {
		return errBadTTL
	}
	if net.ParseIP(a.Value) == nil {
		return errBadValue
	}
	key := a.Name + strconv.Itoa(int(a.Type))
	var val TypeRecord
	val.TTL = a.TTL
	val.Value = a.Value
	mapLock.Lock()
	defer mapLock.Unlock()
	RecordMap[key] = []TypeRecord{val}
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

type recordArgs struct {
	Name  string `json:"name"`
	Type  uint16 `json:"type"`
	TTL   uint32 `json:"ttl"`
	Value string `json:"value"`
}

type recordRet struct {
	Err int    `json:"err"`
	Msg string `json:"msg"`
}

func HandleHTTP() {
	http.HandleFunc("/api/add.record", func(w http.ResponseWriter, r *http.Request) {
		var ret recordRet
		var a recordArgs
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
		return
	})
	http.HandleFunc("/api/list.record", func(w http.ResponseWriter, r *http.Request) {
		l := listRecord()
		bt, _ := json.Marshal(l)
		w.Write(bt)
		return
	})
	go runHTTP()
}
