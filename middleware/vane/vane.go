package vane

import (
	"fmt"
	"net"

	"github.com/miekg/coredns/middleware"
	"github.com/miekg/coredns/middleware/vane/engine"
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
		cip         net.IP
		clientSetID int
		ok          bool
	)

	if len(r.Question) == 0 {
		return 0, nil
	}

	q := r.Question[0]
	if q.Qtype != dns.TypeA {
		return 0, nil
	}

	state := request.Request{W: w, Req: r}

	value := ctx.Value("vane_engine")
	fmt.Println("vane_engine in vane:", value)
	e, ok := value.(*engine.Engine)
	if !ok || e == nil {
		return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
	}

	value = ctx.Value("clientset_id")
	clientSetID, ok = value.(int)
	if !ok {
		cip = state.GetRemoteAddr()
		if cip == nil {
			clientSetID = 1
		} else {
			clientSetID = e.GetClientSetID(cip)
		}

		fmt.Println("self decoede clientset_id", clientSetID)
	} else {
		fmt.Println("fetch clientset_id from ctx:", clientSetID)
	}

	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}
