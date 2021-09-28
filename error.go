package main

import "errors"

// ErrBadQCount is more than 1 question count
var ErrBadQCount = errors.New("bad question count")

// ErrNotSupported only supports A record
var ErrNotSupported = errors.New("type not supported")

var ErrValExists = errors.New("no default line")
var ErrNoSuchVal = errors.New("no such value")
var ErrBadName = errors.New("bad name")
var ErrBadType = errors.New("bad type")
var ErrBadTTL = errors.New("bad ttl")
var ErrBadValue = errors.New("bad value")
var ErrNoGoodServers = errors.New("all servers are dead")
