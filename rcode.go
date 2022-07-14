package main

import (
	"errors"

	"github.com/miekg/dns"
)

var ErrRCode = errors.New("rcode name")
var RcodeMap = make(map[string]int)

// IsRcode 返回0说明没有
func IsRcode(name string) int {
	name = dns.Fqdn(name)
	return RcodeMap[name]
}
