// Copyright 2018 Sean.ZH

// dnslite provides a simple dns service and api

package main

// add AAAA record
// curl -X POST -d '{"name":"dilfish.dev", "type":28, "ttl":100, "value":"2001:470:23:976::2"}' http://127.0.0.1:8085/api/add.record
// add A record
// curl -X POST -d '{"name":"dilfish.dev", "type":1, "ttl":100, "a":"1.1.1.1"}' http://127.0.0.1:8085/api/add.record
// add TXT record
// curl -X POST -d '{"name":"dilfish.dev", "type":16, "ttl":100, "txt":"abcde"}' http://127.0.0.1:8085/api/add.record
// add CNAME record
// curl -X POST -d '{"name":"dilfish.dev", "type":5, "ttl":100, "cname":"abcde.com"}' http://127.0.0.1:8085/api/add.record
// add NS record
// curl -X POST -d '{"name":"www.dilfish.dev", "type":2, "ttl":100, "ns":"abcde.com"}' http://127.0.0.1:8085/api/add.record
// add CAA record
// curl -X POST -d '{"name":"www.dilfish.dev", "type":257, "ttl":100, "caaTag":"issue", "caaFlag":1, "caaValue":"111"}' http://127.0.0.1:8085/api/add.record
// add SVCB record
// curl -X POST -d '{"name":"svcb.dilfish.dev","type":64, "svcbAlpn":{"Alpn":["a","b"]},"svcbIPv6Hint":{"Hint":["1::1"]},"svcbIPv4Hint":{"Hint":["1.1.1.1"]},"ttl":100,"svcbPriority":1,"svcbTarget":"dilfish.dev"}' http://127.0.0.1:8085/api/add.record

// get record list
// curl -X POST -d '{"name":"dilfish.dev","type":1}' http://127.0.0.1:8085/api/list.record

// del a record
// curl -X POST -d '{"name":"dilfish.dev", "type":1}' http://127.0.0.1:8085/api/del.record
