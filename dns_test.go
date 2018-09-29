package dnslite

import (
	"github.com/miekg/dns"
	"net"
	"testing"
)

type FakeResponseWriter struct {
	m    []dns.Msg
	b    []byte
	tsig bool
}

func (f *FakeResponseWriter) LocalAddr() net.Addr {
	a, _ := net.InterfaceAddrs()
	return a[0]
}

func (f *FakeResponseWriter) RemoteAddr() net.Addr {
	a, _ := net.InterfaceAddrs()
	return a[0]
}

func (f *FakeResponseWriter) WriteMsg(m *dns.Msg) error {
	f.m = append(f.m, *m)
	return nil
}

func (f *FakeResponseWriter) Write(b []byte) (int, error) {
	f.b = append(f.b, b...)
	return len(b), nil
}

func (f *FakeResponseWriter) Close() error {
	return nil
}

func (f *FakeResponseWriter) TsigStatus() error {
	return nil
}

func (f *FakeResponseWriter) TsigTimersOnly(set bool) {
	f.tsig = set
}

func (f *FakeResponseWriter) Hijack() {}

func TestCreateDNSMux(t *testing.T) {
	mux := CreateDNSMux()
	var m dns.Msg
	var w FakeResponseWriter
	m.Id = dns.Id()
	m.RecursionDesired = true
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{"sub.ns.libsm.com.", dns.TypeA, dns.ClassINET}
	mux.ServeDNS(&w, &m)
	if len(w.b) == 0 && len(w.m) == 0 {
		t.Error("return nil")
	}
}
