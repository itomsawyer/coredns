package vane

import (
	"fmt"
	"net"

	"github.com/miekg/coredns/middleware"
	"github.com/miekg/coredns/middleware/pkg/edns"
	"github.com/miekg/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type Vane struct {
	Next   middleware.Handler
	DB     Loader
	engine *Engine
	DBHost string
}

func (v Vane) Name() string { return "vane" }

func (v Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		remoteAddr net.IP
		engine     *Engine
	)

	engine = v.engine

	if len(r.Question) == 0 {
		return 0, nil
	}

	q := r.Question[0]
	if q.Qtype != dns.TypeA {
		return 0, nil
	}

	state := request.Request{W: w, Req: r}

	subnet := edns.ReadClientSubnet(r)
	if subnet == nil {
		remoteAddr = net.ParseIP(state.IP()).To4()
	} else {
		remoteAddr = subnet.Address
	}

	clientSet, err := engine.GetClientSetID(remoteAddr)
	if err != nil {
		clientSet = 1
	}

	fmt.Println("Client Set:", clientSet)

	fmt.Println("Proxys:", engine.Proxy)
	fmt.Println("Exchange to", engine.Proxy[1])
	_, err = engine.Proxy[1].Lookup(state, q.Name, q.Qtype)
	if err != nil {
		fmt.Println("error ", err)
		return 0, err
	}
	fmt.Println("done")

	answer := new(dns.Msg)
	answer.SetReply(r)
	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
	rr.A = net.ParseIP("127.0.0.1")
	answer.Answer = []dns.RR{rr}

	w.WriteMsg(answer)

	return 0, nil
}
