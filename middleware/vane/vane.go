package vane

import (
	"net"

	"github.com/miekg/coredns/middleware"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type Vane struct {
	Next middleware.Handler
}

func (d Vane) Name() string { return "vane" }

func (d Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		return 0, nil
	}

	q := r.Question[0]
	if q.Qtype != dns.TypeA {
		return 0, nil
	}

	answer := new(dns.Msg)
	answer.SetReply(r)
	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
	rr.A = net.ParseIP("127.0.0.1")
	answer.Answer = []dns.RR{rr}

	w.WriteMsg(answer)

	return 0, nil
}
