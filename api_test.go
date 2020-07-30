// Copyright 2018 Sean.ZH

package dnslite

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

type DNSApi struct {
	t *testing.T
}

func (d *DNSApi) Add(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	data := []byte(r.Body.String())
	jerr, _ := jsonparser.GetInt(data, "err")
	jmsg, _ := jsonparser.GetString(data, "msg")
	assert.Equal(d.t, int64(0), jerr)
	assert.Equal(d.t, "ok", jmsg)
	assert.Equal(d.t, http.StatusOK, r.Code)
}

func (d *DNSApi) List(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	var ris []RecordInfo
	err := json.Unmarshal([]byte(r.Body.String()), &ris)
	if err != nil {
		d.t.Error("unjson", err)
		return
	}
	if len(ris) != 1 {
		d.t.Error("ris is", len(ris))
		return
	}
	if ris[0].Name != "sub.ns.libsm.com.1" || ris[0].Value != "1.1.1.1" || ris[0].TTL != 5 {
		d.t.Error("values error", ris[0])
	}
}

func TestCreateHTTPMux(t *testing.T) {
	mux := CreateHTTPMux()
	r := gofight.New()
	var d DNSApi
	d.t = t
	r.POST("/api/add.record").
		SetJSON(gofight.D{
			"name":  "sub.ns.libsm.com.",
			"type":  1,
			"ttl":   5,
			"value": "1.1.1.1",
		}).SetDebug(true).Run(mux, d.Add)
	r.POST("/api/list.record").SetDebug(true).Run(mux, d.List)
}
