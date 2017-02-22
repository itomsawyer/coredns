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
		replyMsg    *dns.Msg
	)

	if len(r.Question) == 0 {
		return 0, nil
	}

	q := r.Question[0]
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

try_again:

	view, err := e.GetView(clientSetID, dmPoolID)
	if err != nil {
		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}
	fmt.Printf("get view %+v \n", view)

	if view.Upstream == nil {
		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}

	for _, e := range view.Upstream.Hosts {
		fmt.Println("host in upstream", e)
	}

	policy := view.Upstream.GetPolicy()
	if policy == nil {

		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}

	fmt.Println("get policy success")

	for {
		fmt.Println("begin select upstream hosts")
		uphosts := policy.Select()
		if len(uphosts) == 0 {
			// There no upstream host can be found , try again the whole precedure with domainPoolID equals to
			// default domainPoolID which is 1. Or the domainPoolID is already the default, Lookup failed the
			// WARNING MUST be sent

			if dmPoolID != engine.DefaultDomainPoolID {
				dmPoolID = engine.DefaultDomainPoolID
				goto try_again
			}

			return dns.RcodeServerFailure, errUnreachable
		}

		for _, uh := range uphosts {
			fmt.Println("found upstream host", uh)
		}

		replys := DNSExWithTimeout(uphosts, state, 1*time.Second)
		for _, r := range replys {
			fmt.Println("get replys", r)
		}

		if len(replys) == 0 {
			fmt.Println("get no replys")
			continue
		}

		if q.Qtype != dns.TypeA {
			if len(replys) > 0 {
				w.WriteMsg(replys[0])
				return 0, nil
			}
			continue
		}

		better := addrSet{}
		for _, reply := range replys {
			rrset := reply.Answer
			for _, rr := range rrset {
				if a, ok := rr.(*dns.A); ok {
					netLinkID := e.GetNetLinkID(a.A)
					routes := e.GetRoute(view.RouteSetID, dmPoolID, netLinkID)
					if len(routes) > 0 {
						if replyMsg == nil {
							replyMsg = reply
						}
						better.Add(a)
					}
				}
			}
		}

		if len(better) > 0 {
			replyMsg.Answer = better.RRSet()
			fmt.Println("write anwser:", replyMsg)
			w.WriteMsg(replyMsg)
			return 0, nil
		}
	}

	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}

func DNSExWithTimeout(uphosts []*proxy.UpstreamHost, state request.Request, timeout time.Duration) (replys []*dns.Msg) {
	if len(uphosts) == 1 {
		uh := uphosts[0]

		reply, backendErr := uh.DoExchange(state)
		if backendErr != nil {
			//TODO LOG
			fmt.Println(backendErr)
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

type addrSet map[string]*dns.A

func (p addrSet) Add(a *dns.A) {
	if a == nil || a.A == nil {
		return
	}

	p[a.A.String()] = a
}

func (p addrSet) RRSet() []dns.RR {
	if len(p) == 0 {
		return nil
	}

	s := make([]dns.RR, 0, len(p))
	for _, a := range p {
		s = append(s, a)
	}

	return s
}
