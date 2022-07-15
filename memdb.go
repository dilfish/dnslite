package main

import (
	"log"
	"sync"
)

func NewDB(conf *Config) DataManagerI {
	if conf.UsingMemDB {
		return NewMemDB()
	}
	if *FlagUsingMongo == "" {
		log.Println("using mongodb but no addr provided")
		return nil
	}
	return NewMongoClient(&conf.MongoClientConfig)
}

type MemDB struct {
	lock sync.Mutex
	mp   map[string]map[uint16][]DNSRecord
}

func NewMemDB() *MemDB {
	var m MemDB
	m.mp = make(map[string]map[uint16][]DNSRecord)
	return &m
}

func (m *MemDB) Find(name string, tp uint16) (ret []DNSRecord, err error) {
	log.Println("memdb find:", name, tp)
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.mp[name]
	if !ok {
		return nil, nil
	}
	vv, ok := v[tp]
	if !ok {
		return nil, nil
	}
	log.Println("memdb found:", name, tp, vv)
	return vv, nil
}

func (m *MemDB) Insert(r DNSRecord) error {
	log.Println("memdb insert:", r.Name, r.Type, r)
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.mp[r.Name]
	if !ok {
		v = make(map[uint16][]DNSRecord)
	}
	vv := v[r.Type]
	vv = append(vv, r)
	v[r.Type] = vv
	m.mp[r.Name] = v
	return nil
}

func (m *MemDB) Del(name string, tp uint16) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.mp[name]
	if !ok {
		return nil
	}
	_, ok = v[tp]
	if !ok {
		return nil
	}
	delete(v, tp)
	m.mp[name] = v
	return nil
}
