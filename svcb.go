// sean at shanghai 2021

package main

import (
	"log"

	"github.com/miekg/dns"
)

type SVCBHandler struct{}

func (svcb *SVCBHandler) FillRecords(req *dns.Msg, records []DNSRecord) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(req)
	rr := make([]dns.SVCB, len(records))
	for idx, record := range records {
		rr[idx].Hdr.Name = req.Question[0].Name
		rr[idx].Hdr.Rrtype = dns.TypeSVCB
		rr[idx].Hdr.Class = dns.ClassINET
		rr[idx].Hdr.Ttl = record.Ttl
		rr[idx].Target = record.SVCBTarget
		rr[idx].Priority = record.SVCBPriority
		if len(record.SVCBIPv4Hint.Hint) != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBIPv4Hint)
		}
		if len(record.SVCBIPv6Hint.Hint) != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBIPv6Hint)
		}
		if len(record.SVCBAlpn.Alpn) != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBAlpn)
		}
		if record.SVCBPort.Port != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBPort)
		}
		if len(record.SVCBMandatory.Code) != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBMandatory)
		}
		if len(record.SVCBECHConfig.ECH) != 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBECHConfig)
		}
		if len(rr[idx].Value) == 0 {
			rr[idx].Value = append(rr[idx].Value, record.SVCBNoDefaultAlpn)
		}
		m.Answer = append(m.Answer, &rr[idx])
	}
	return m
}

func (svcb *SVCBHandler) CheckRecord(record *DNSRecord) error {
	if record.SVCBPriority == 0 {
		log.Println("bad svcb priority")
		return ErrBadValue
	}
	if len(record.SVCBTarget) == 0 {
		log.Println("bad svcb target")
		return ErrBadValue
	}
	record.SVCBTarget = AppendDot(record.SVCBTarget)
	return nil
}

func (svcb *SVCBHandler) RRToRecord(msg dns.RR) DNSRecord {
	var record DNSRecord
	s := msg.(*dns.SVCB)
	record.SVCBPriority = s.Priority
	record.SVCBTarget = s.Target
	for _, v := range s.Value {
		switch v.Key() {
		case dns.SVCB_IPV6HINT:
			vv := v.(*dns.SVCBIPv6Hint)
			record.SVCBIPv6Hint = vv
		case dns.SVCB_IPV4HINT:
			vv := v.(*dns.SVCBIPv4Hint)
			record.SVCBIPv4Hint = vv
		case dns.SVCB_ALPN:
			vv := v.(*dns.SVCBAlpn)
			record.SVCBAlpn = vv
		case dns.SVCB_PORT:
			vv := v.(*dns.SVCBPort)
			record.SVCBPort = vv
		case dns.SVCB_MANDATORY:
			vv := v.(*dns.SVCBMandatory)
			record.SVCBMandatory = vv
		case dns.SVCB_ECHCONFIG:
			vv := v.(*dns.SVCBECHConfig)
			record.SVCBECHConfig = vv
		}
	}
	return record
}
