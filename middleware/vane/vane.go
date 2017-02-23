package vane

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var (
	errUnreachable = errors.New("unreachable backend")
	errFormatError = errors.New("format error")
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
		return dns.RcodeFormatError, errFormatError
	}

	q := r.Question[0]
	state := request.Request{W: w, Req: r}

	//XXX MayBe check vane_engine is start at first startup
	value := ctx.Value("vane_engine")
	e, ok := value.(*engine.Engine)
	if !ok || e == nil {
		return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
	}

	// Try get clientset_id from previous vane_engine middleware which has done this job.
	// In case vane_engine doesn't do its duty, try here then
	value = ctx.Value("clientset_id")
	clientSetID, ok = value.(int)
	if !ok {
		cip = state.GetRemoteAddr()
		if cip == nil {
			clientSetID = engine.DefaultClientSetID
		} else {
			clientSetID = e.GetClientSetID(cip)
		}

		fmt.Println("self decoede clientset_id", clientSetID)
	} else {
		fmt.Println("fetch clientset_id from ctx:", clientSetID)
	}

	// Get domainpool_id , if not found e.GetDomainPoolID return engine.DefaultDomainPoolID
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

	//Policy is a method class to choose upstreamhost (ldns) from upstream (policy)
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
		// for each time policy choose the next prior upstreamhosts group
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

		// Send dns query to every upstreamhost in uphosts, combine their response into slice replys
		// TODO make the timeout configurable
		replys := DNSExWithTimeout(ctx, view.Upstream, uphosts, state, 1*time.Second)
		for _, r := range replys {
			fmt.Println("get replys", r)
		}

		if len(replys) == 0 {
			fmt.Println("get no replys")
			continue
		}

		// No need to filter record with type is not, Get a proper one to return
		if q.Qtype != dns.TypeA {
			if len(replys) > 0 {
				w.WriteMsg(replys[0])
				return 0, nil
			}
			continue
		}

		// better is the result set of all A that pass the filter with Route
		better := addrSet{}
		for _, reply := range replys {
			rrset := reply.Answer
			for _, rr := range rrset {
				if a, ok := rr.(*dns.A); ok {
					netLinkID := e.GetNetLinkID(a.A)
					routes := e.GetRoute(view.RouteSetID, dmPoolID, netLinkID)
					// If has route, we consider the result to be valid
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
			// we got answer, return
			replyMsg.Answer = better.RRSet()
			fmt.Println("write anwser:", replyMsg)
			w.WriteMsg(replyMsg)
			return 0, nil
		}

		// No luck for this time, try to ask other upstreamhosts
	}

	// TODO LOG WARN: we tried our best but still got nothing
	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}

func DNSExWithTimeout(ctx context.Context, upstream *engine.Upstream, uphosts []*proxy.UpstreamHost, state request.Request, timeout time.Duration) (replys []*dns.Msg) {
	ex := upstream.Exchanger()
	if ex == nil {
		//TODO LOG
		return nil
	}

	if len(uphosts) == 1 {
		uh := uphosts[0]

		atomic.AddInt64(&uh.Conns, 1)
		reply, backendErr := ex.Exchange(ctx, uh.Name, state)
		atomic.AddInt64(&uh.Conns, -1)

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
	done := make(chan struct{})

	for i := 0; i < len(uphosts); i++ {
		wg.Add(1)
		go func(uh *proxy.UpstreamHost) {
			fmt.Println("send query to ", uh)
			atomic.AddInt64(&uh.Conns, 1)
			reply, backendErr := ex.Exchange(ctx, uh.Name, state)
			atomic.AddInt64(&uh.Conns, -1)
			fmt.Println("exchange get reply with error:", backendErr, "msg:", reply)

			if backendErr == nil {
				select {
				case out <- reply:
				case <-done:
				}
			} else {
				select {
				case errChan <- backendErr:
				case <-done:
				}
				uh.Fail()
			}

			wg.Done()
		}(uphosts[i])
	}

	defer func() {
		close(done)
		go func() {
			fmt.Println("clean up DNSExWithTimeout: waiting")
			wg.Wait()
			close(out)
			close(errChan)
			fmt.Println("clean up DNSExWithTimeout: done")
		}()
	}()

	if timeout == 0 {
		// TODO make it configurable
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
