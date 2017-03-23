package vane

import (
	"errors"
	"math"
	"net"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"

	"github.com/astaxie/beego/logs"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

var (
	errUnreachable     = errors.New("unreachable backend")
	errUnexpectedLogic = errors.New("incorrect or incomplete data logic")
	errFormatError     = errors.New("format error")

	defaultMaxAttampts = 10
	defaultTimeout     = 3 * time.Second
)

type Vane struct {
	UpstreamTimeout time.Duration
	Debug           bool
	Logger          *logs.BeeLogger
	LogConfigs      []*engine.LogConfig
	RcodePriority   *RcodePriority
	Next            middleware.Handler
}

func NewVane() *Vane {
	return &Vane{
		RcodePriority: NewRcodePriority(),
	}
}

func (v *Vane) Init() error {
	l, err := engine.CreateLogger(v.LogConfigs)
	if err != nil {
		return err
	}

	v.Logger = l
	v.Logger.Info("vane start success")
	return nil
}

func (v *Vane) Destroy() error {
	if v.Logger != nil {
		v.Logger.Close()
	}

	return nil
}

func (v *Vane) Name() string { return "vane" }

func (v *Vane) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		cip             net.IP
		i, clientSetID  int
		ok              bool
		replyMsg        *dns.Msg
		replys          []*dns.Msg
		retcode, bestrc int
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
	cip = state.GetRemoteAddr()
	v.Logger.Debug("get remote client addr %v", cip)
	clientSetID, ok = value.(int)
	if !ok {
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
	domain, err := e.GetDomain(q.Name)
	if err != nil {
		domain = engine.DefaultDomainPool
	}

	v.Logger.Debug("query domain %s", q.Name)
	v.Logger.Debug("get domain pool id: %d", domain.DmPoolID)

try_again:

	view, err := e.GetView(clientSetID, domain.DmPoolID)
	if err != nil {
		if clientSetID, domain, ok = rollback(clientSetID, domain); ok {
			goto try_again
		}

		return dns.RcodeServerFailure, errUnexpectedLogic
	}

	v.Logger.Debug("get view %v", view)

	if view.Upstream == nil {
		if clientSetID, domain, ok = rollback(clientSetID, domain); ok {
			goto try_again
		}

		return dns.RcodeServerFailure, errUnexpectedLogic
	}

	//Policy is a method class to choose upstreamhost (ldns) from upstream (policy)
	policy := view.Upstream.GetPolicy()
	if policy == nil {
		if clientSetID, domain, ok = rollback(clientSetID, domain); ok {
			goto try_again
		}

		return dns.RcodeServerFailure, errUnexpectedLogic
	}

	ex := NewExchangeHelper(view.Upstream, nil)
	ex.Timeout = v.UpstreamTimeout
	bestrc = dns.RcodeServerFailure
	replyMsg = nil
	replys = nil

	for i = 0; i < defaultMaxAttampts; i++ {
		v.Logger.Debug("try select upstream hosts")
		// for each time policy choose the next prior upstreamhosts group
		uphosts := policy.Select()
		if len(uphosts) == 0 {
			// There no upstream host can be found , check out what we got.
			break
		}

		if v.Debug {
			for _, uh := range uphosts {
				v.Logger.Debug("upstream host found %v", uh)
			}
		}

		// Send dns query to every upstreamhost in uphosts, combine their response into slice replys
		v.Logger.Debug("set upstream timeout: %s", v.UpstreamTimeout)
		ex.Hosts = uphosts
		replys, retcode = ex.DoExchange(ctx, state)
		if v.Debug {
			for _, r := range replys {
				v.Logger.Debug("get reply from upstream host :\n%v", r)
			}
		}

		if len(replys) == 0 {
			continue
		}

		v.Logger.Debug("found retcode %d with bestrc cmp to current %d", retcode, bestrc)
		if v.RcodePriority.PriorTo(retcode, bestrc) {
			v.Logger.Debug("found better retcode: %d", retcode)
			v.Logger.Debug("it says: \n%v", replys)
			bestrc = retcode
			replyMsg = replys[0]
		}

		if retcode != dns.RcodeSuccess {
			continue
		}

		// No need to filter record with type is not A, Get a proper one to return
		if q.Qtype != dns.TypeA {
			if bestrc == dns.RcodeSuccess && replyMsg != nil {
				w.WriteMsg(replyMsg)
				return 0, nil
			}

			continue
		}

		// better is the result set of all A that pass the filter with Route
		better := make([]dns.RR, 0, 4)
		rrset := rrSet{}
		for _, reply := range replys {
			rrlist := reply.Answer
			for _, rr := range rrlist {
				if a, ok := rr.(*dns.A); ok {
					rrset.Add(a)
				}
			}
		}

		bestRoutePrio := math.MaxInt32
		for _, rr := range rrset {
			a := rr.(*dns.A)
			netLinkID := e.GetNetLinkID(a.A)
			v.Logger.Debug("ip %s found netLinkID: %d", a.A, netLinkID)
			routes := e.GetRoute(view.RouteSetID, domain.DmPoolID, netLinkID)
			// If has route, we consider the result to be valid
			if len(routes) == 0 && netLinkID != engine.DefaultNetLinkID {
				v.Logger.Debug("no route found try default netlinkID instead")
				routes = e.GetRoute(view.RouteSetID, domain.DmPoolID, engine.DefaultNetLinkID)
			}

			if len(routes) == 0 {
				v.Logger.Debug("no route found", netLinkID)
				continue
			}

			route := routes[0]

			if domain.Monitor {
				v.Logger.Debug("domain %s with dmpool %d need to use dynamic monitor", q.Name, domain.DmPoolID)
				status, ok := e.LinkManager.GetLink(a.String(), routes[0].OutLink.Addr)
				if ok && status.Status == engine.LinkStatusDown {
					v.Logger.Debug("%s dynamic monitor status down", a)
					continue
				}
			}

			if route.Priority < bestRoutePrio {
				bestRoutePrio = route.Priority
				better = better[:0]
			}

			if route.Priority > bestRoutePrio {
				continue
			}

			v.Logger.Debug("ip addr %s has been accepted for %s", a.A, route.OutLink.Addr)
			a.Hdr.Name = q.Name
			better = append(better, a)
		}

		if len(better) == 0 {
			// No luck for this time, try to ask other upstreamhosts
			continue
		}

		// we got answer, return
		v.Logger.Debug("Rewrite Msg: %v\n", replyMsg)
		replyMsg.Answer = better
		v.Logger.Debug("Write anwser to client: \n%v", replyMsg)
		w.WriteMsg(replyMsg)
		return 0, nil
	}

	if i == defaultMaxAttampts {
		v.Logger.Error("MaxAttampts of upstream query reached, configuration or policy maybe badly configured")
	}

	// try again the whole precedure with domainPoolID equals to
	// default domainPoolID which is 1. Or the domainPoolID is already the default, Lookup failed

	if clientSetID, domain, ok = rollback(clientSetID, domain); ok {
		goto try_again
	}

	// TODO WARNING MUST be sent when we get here
	// 1. replyMsg is not nil , bestrc is dns.RcodeSuccess: we filtered all of ip addr. At least, give out one answer
	// 2. domain pool has falled back to default once and still got no good answer

	//return the best effort answer we get
	if replyMsg != nil {
		v.Logger.Warn("vane cannot find a good answer for client: %v domain: %s", cip, q.Name)
		w.WriteMsg(replyMsg)
		return 0, nil
	}

	//we tried our best but still got nothing
	v.Logger.Error("vane cannot resolve any answer for client: %v domain: %s", cip, q.Name)
	return dns.RcodeServerFailure, errUnreachable
}

func rollback(clientSetID int, domain engine.Domain) (int, engine.Domain, bool) {
	if domain.DmPoolID != engine.DefaultDomainPool.DmPoolID {
		return clientSetID, engine.DefaultDomainPool, true
	}

	if clientSetID != engine.DefaultClientSetID {
		return engine.DefaultClientSetID, domain, true
	}

	return 0, engine.Domain{}, false
}
