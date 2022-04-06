package main

import (
	"sync"
)

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
	return vv, nil
}

func (m *MemDB) Insert(r DNSRecord) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	v, ok := m.mp[r.Name]
	if !ok {
		v = make(map[uint16][]DNSRecord)
	}
	vv := v[r.Type]
	vv = append(vv, r)
	v[r.Type] = vv
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
