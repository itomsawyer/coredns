package vane

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"

	"github.com/astaxie/beego/logs"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var (
	errUnreachable = errors.New("unreachable backend")
	errFormatError = errors.New("format error")

	defaultMaxAttampts = 10
	defaultTimeout     = 3 * time.Second
)

type Vane struct {
	UpstreamTimeout time.Duration
	Debug           bool
	Logger          *logs.BeeLogger
	LogConfigs      []*engine.LogConfig
	Next            middleware.Handler
}

func (v *Vane) InitLogger() error {
	l, err := engine.CreateLogger(v.LogConfigs)
	if err != nil {
		return err
	}

	v.Logger = l
	return nil
}

func (v *Vane) Name() string { return "vane" }

func (v *Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		cip            net.IP
		i, clientSetID int
		ok             bool
		replyMsg       *dns.Msg
	)

	if len(r.Question) == 0 {
		return dns.RcodeFormatError, errFormatError
	}

	q := r.Question[0]
	state := request.Request{W: w, Req: r}

	//TODO MayBe check vane_engine is start at first startup
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
		v.Logger.Debug("get remote client addr %v", cip)
		if cip == nil {
			clientSetID = engine.DefaultClientSetID
		} else {
			clientSetID = e.GetClientSetID(cip)
		}

		v.Logger.Debug("fetch clientset_id from request: %v", clientSetID)
	} else {
		v.Logger.Debug("fetch clientset_id from vane_engine: %v", clientSetID)
	}

	// Get domainpool_id , if not found e.GetDomainPoolID return engine.DefaultDomainPoolID
	dmPoolID := e.GetDomainPoolID(q.Name)
	v.Logger.Debug("query domain %s", q.Name)
	v.Logger.Debug("get domain pool id: %d", dmPoolID)

try_again:

	view, err := e.GetView(clientSetID, dmPoolID)
	if err != nil {
		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			v.Logger.Debug("fallback to default domain pool")
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}

	v.Logger.Debug("get view %v", view)

	if view.Upstream == nil {
		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			v.Logger.Debug("fallback to default domain pool")
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}

	//Policy is a method class to choose upstreamhost (ldns) from upstream (policy)
	policy := view.Upstream.GetPolicy()
	if policy == nil {

		if dmPoolID != engine.DefaultDomainPoolID {
			dmPoolID = engine.DefaultDomainPoolID
			v.Logger.Debug("fallback to default domain pool")
			goto try_again
		}

		return dns.RcodeServerFailure, errUnreachable
	}

	for i = 0; i < defaultMaxAttampts; i++ {
		v.Logger.Debug("try select upstream hosts")
		// for each time policy choose the next prior upstreamhosts group
		uphosts := policy.Select()
		if len(uphosts) == 0 {
			// There no upstream host can be found , try again the whole precedure with domainPoolID equals to
			// default domainPoolID which is 1. Or the domainPoolID is already the default, Lookup failed the
			// WARNING MUST be sent

			if dmPoolID != engine.DefaultDomainPoolID {
				dmPoolID = engine.DefaultDomainPoolID
				v.Logger.Debug("fallback to default domain pool")
				goto try_again
			}

			return dns.RcodeServerFailure, errUnreachable
		}

		if v.Debug {
			for _, uh := range uphosts {
				v.Logger.Debug("upstream host found %v", uh)
			}
		}

		// Send dns query to every upstreamhost in uphosts, combine their response into slice replys
		v.Logger.Debug("set upstream timeout: %s", v.UpstreamTimeout)
		replys := DNSExWithTimeout(ctx, view.Upstream, uphosts, state, v.UpstreamTimeout)
		if v.Debug {
			for _, r := range replys {
				v.Logger.Debug("get reply from upstream host :\n%v", r)
			}
		}

		if len(replys) == 0 {
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
						v.Logger.Debug("ip addr has been accepted %s", a.A)
						better.Add(a)
					}
				}
			}
		}

		if len(better) > 0 {
			// we got answer, return
			replyMsg.Answer = better.RRSet()
			v.Logger.Debug("Write anwser to client: \n%v", replyMsg)
			w.WriteMsg(replyMsg)
			return 0, nil
		}

		// No luck for this time, try to ask other upstreamhosts
	}

	if i == defaultMaxAttampts {
		v.Logger.Error("MaxAttampts of upstream query reached, configuration or policy is badly configured")
	}

	if dmPoolID != engine.DefaultDomainPoolID {
		dmPoolID = engine.DefaultDomainPoolID
		v.Logger.Warn("fallback to default domain pool for clientset: %d dmpool: %d", clientSetID, dmPoolID)
		goto try_again
	}

	//LOG WARN: we tried our best but still got nothing
	v.Logger.Error("vane cannot resolve for clientset: %d dmpool: %d", clientSetID, dmPoolID)

	return middleware.NextOrFailure(v.Name(), v.Next, ctx, w, r)
}

func DNSExWithTimeout(ctx context.Context, upstream *engine.Upstream, uphosts []*proxy.UpstreamHost, state request.Request, timeout time.Duration) (replys []*dns.Msg) {
	ex := upstream.Exchanger()
	if ex == nil {
		return nil
	}

	if len(uphosts) == 1 {
		uh := uphosts[0]

		atomic.AddInt64(&uh.Conns, 1)
		reply, backendErr := ex.Exchange(ctx, uh.Name, state)
		atomic.AddInt64(&uh.Conns, -1)

		if backendErr != nil {
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
			atomic.AddInt64(&uh.Conns, 1)
			reply, backendErr := ex.Exchange(ctx, uh.Name, state)
			atomic.AddInt64(&uh.Conns, -1)

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
			wg.Wait()
			close(out)
			close(errChan)
		}()
	}()

	if timeout == 0 {
		timeout = defaultTimeout
	}

	t := time.NewTimer(timeout)

	for cnt := 0; cnt < len(uphosts); cnt++ {
		select {
		case reply := <-out:
			replys = append(replys, reply)
		case <-errChan:
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
