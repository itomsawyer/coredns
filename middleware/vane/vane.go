package vane

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/miekg/coredns/middleware"
	"github.com/miekg/coredns/middleware/proxy"
	"github.com/miekg/coredns/middleware/vane/engine"
	"github.com/miekg/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var (
	errUnreachable = errors.New("unreachable backend")
)

type Vane struct {
	Next middleware.Handler
}

func (v Vane) Name() string { return "vane" }

func (v Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	fmt.Println("enter vane")
	defer fmt.Println("leave vane")
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

	dmPoolID := e.GetDomainPoolID(q.Name)
	fmt.Println("get domain pool id", dmPoolID)

	view, err := e.GetView(clientSetID, dmPoolID)
	if err != nil {
		return dns.RcodeServerFailure, errUnreachable
	}
	fmt.Printf("get view %+v \n", view)

	if view.Upstream == nil {
		return dns.RcodeServerFailure, errUnreachable
	}

	policy := view.Upstream.GetPolicy()
	if policy == nil {
		return dns.RcodeServerFailure, errUnreachable
	}

	fmt.Println("get policy success")

	for {
		fmt.Println("begin select upstream hosts")
		uphosts := policy.Select()
		if len(uphosts) == 0 {
			//TODO goto try_again
			break
		}

		fmt.Println("found upstream host", uphosts)

		replys := DNSExWithTimeout(uphosts, state, 1*time.Second)

		better := make([]dns.RR, 0, 1)
		for _, reply := range replys {
			rrset := reply.Answer
			for _, rr := range rrset {
				if a, ok := rr.(*dns.A); ok {
					netLinkID := e.GetNetLinkID(a.A)
					routes := e.GetRoute(view.RouteSetID, dmPoolID, netLinkID)
					if len(routes) > 0 {
						better = append(better, rr)
					}
				}
			}
		}

		if len(better) == 0 {

		}
	}

	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}

func DNSExWithTimeout(uphosts []*proxy.UpstreamHost, state request.Request, timeout time.Duration) (replys []*dns.Msg) {
	if len(uphosts) == 1 {
		uh := uphosts[0]

		reply, backendErr := uh.DoExchange(state)
		if backendErr == nil {
			//TODO LOG
			uh.Fail()
			return nil
		}

		return []*dns.Msg{reply}
	}

	out := make(chan *dns.Msg)
	errChan := make(chan error)
	wg := new(sync.WaitGroup)

	for i := 0; i < len(uphosts); i++ {
		wg.Add(1)
		go func(uh *proxy.UpstreamHost) {
			reply, backendErr := uh.DoExchange(state)

			if backendErr == nil {
				out <- reply
			} else {
				errChan <- backendErr
				uh.Fail()
			}

			wg.Done()
		}(uphosts[i])
	}

	defer func() {
		go func() {
			wg.Wait()
			close(out)
			close(errChan)
		}()
	}()

	if timeout == 0 {
		timeout = 3 * time.Second
	}

	t := time.NewTimer(timeout)

	for cnt := 0; cnt < len(uphosts); cnt++ {
		select {
		case reply := <-out:
			replys = append(replys, reply)
		case <-errChan:
			//TODO LOG
		case <-t.C:
			return
		}
	}

	return
}
