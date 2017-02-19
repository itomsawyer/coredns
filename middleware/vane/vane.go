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
	Next middleware.Handler
}

func (v Vane) Name() string { return "vane" }

func (v Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		remoteAddr net.IP
	)

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

	fmt.Println(remoteAddr)

	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}
