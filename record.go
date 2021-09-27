package dnslite

import (
	"github.com/miekg/dns"
)

type DNSRecord struct {
	Id                string                `json:"id" bson:"id"`
	Name              string                `json:"name" bson:"name"`
	Type              uint16                `json:"type" bson:"type"`
	Ttl               uint32                `json:"ttl" bson:"ttl"`
	View              string                `json:"view" bson:"view"`
	A                 string                `json:"a,omitempty" bson:"a,omitempty"`
	Aaaa              string                `json:"aaaa,omitempty" bson:"aaaa,omitempty"`
	Ns                string                `json:"ns,omitempty" bson:"ns,omitempty"`
	Cname             string                `json:"cname,omitempty" bson:"cname,omitempty"`
	Txt               string                `json:"txt,omitempty" bson:"txt,omitempty"`
	CAATag            string                `json:"caaTag,omitempty" bson:"caaTag,omitempty"`
	CAAFlag           uint8                 `json:"CaaFlag,omitempty" bson:"caaFlag,omitempty"`
	CAAValue          string                `json:"caaValue,omitempty" bson:"caaValue,omitempty"`
	Code              int                   `json:"code,omitempty" bson:"code,omitempty"`
	Msg               string                `json:"msg,omitempty" bson:"msg,omitempty"`
	SVCBTarget        string                `json:"svcbTarget,omitempty" bson:"svcbTarget,omitempty"`
	SVCBPriority      uint16                `json:"svcbPriority,omitempty" bson:"svcbPriority,omitempty"`
	SVCBPort          dns.SVCBPort          `json:"svcbPort,omitempty" bson:"svcbPort,omitempty"`
	SVCBMandatory     dns.SVCBMandatory     `json:"svcbMandatory,omitempty" bson:"svcbMandatory,omitempty"`
	SVCBAlpn          dns.SVCBAlpn          `json:"svcbAlpn,omitempty" bson:"svcbAlpn,omitempty"`
	SVCBECHConfig     dns.SVCBECHConfig     `json:"svcbECHConfig,omitempty" bson:"svcbECHConfig,omitempty"`
	SVCBIPv4Hint      dns.SVCBIPv4Hint      `json:"svcbIPv4Hint,omitempty" bson:"svcbIPv4Hint,omitempty"`
	SVCBIPv6Hint      dns.SVCBIPv6Hint      `json:"svcbIPv6Hint,omitempty" bson:"svcbIPv6Hint,omitempty"`
	SVCBNoDefaultAlpn dns.SVCBNoDefaultAlpn `json:"svcbNoDefaultAlpn,omitempty" bson:"svcbNoDefaultAlpn,omitempty"`
}
