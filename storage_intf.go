package main

type DataManagerI interface {
	Find(name string, tp uint16) (ret []DNSRecord, err error)
	Insert(r DNSRecord) error
	Del(name string, tp uint16) error
}
