package dnslite

import (
	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type DNSApi struct {
	t *testing.T
}

func (d *DNSApi) Add(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
	data := []byte(r.Body.String())
	jerr, _ := jsonparser.GetInt(data, "err")
	jmsg, _ := jsonparser.GetString(data, "msg")
	assert.Equal(d.t, 0, jerr)
	assert.Equal(d.t, "ok", jmsg)
	assert.Equal(d.t, http.StatusOK, r.Code)
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
}
