// Copyright 2018 Sean.ZH

// dnslite provides a simple dns service and api

package dnslite

// add AAAA record
// curl -X POST -d '{"name":"ipv6.xn--e1t.co", "type":28, "ttl":100, "value":"2001:470:23:976::2"}' http://127.0.0.1:8085/api/add.record
// get record list
// curl http://127.0.0.1:8085/api/list.record
// del aaaa record
// curl -X POST -d '{"name":"ipv6.xn--e1t.co", "type":28, "value":"2001:470:23:976::2"}' http://127.0.0.1:8085/api/del.record
