package main

import (
	"log"

	"github.com/miekg/dns"
)

type SoaHandler struct{}

func (soa *SoaHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.SOA, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeSVCB
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Ns = record.SoaNs
		rr[idx].Mbox = record.SoaMbox
		rr[idx].Serial = record.SoaSerial
		rr[idx].Refresh = record.SoaRefresh
		rr[idx].Retry = record.SoaRetry
		rr[idx].Expire = record.SoaExpire
		rr[idx].Minttl = record.SoaMinttl
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (soa *SoaHandler) CheckRecord(record *DNSRecord) error {
	if len(record.SoaNs) == 0 {
		log.Println("bad soa ns")
		return ErrBadValue
	}
	if len(record.SoaMbox) == 0 {
		log.Println("bad soa mbox")
		return ErrBadValue
	}
	if record.SoaSerial == 0 {
		log.Println("bad soa serial")
		return ErrBadValue
	}
	if record.SoaRefresh == 0 {
		log.Println("bad soa refresh")
		return ErrBadValue
	}
	if record.SoaRetry == 0 {
		log.Println("bad soa retry")
		return ErrBadValue
	}
	if record.SoaExpire == 0 {
		log.Println("bad soa expire")
		return ErrBadValue
	}
	if record.SoaMinttl == 0 {
		log.Println("bad soa minttl")
		return ErrBadValue
	}
	return nil
}

func (soa *SoaHandler) RRToRecord(msg dns.RR) DNSRecord {
	var record DNSRecord
	v := msg.(*dns.SOA)
	record.SoaNs = v.Ns
	record.SoaMbox = v.Mbox
	record.SoaSerial = v.Serial
	record.SoaRefresh = v.Refresh
	record.SoaRetry = v.Retry
	record.SoaExpire = v.Expire
	record.SoaMinttl = v.Minttl
	return record
}
