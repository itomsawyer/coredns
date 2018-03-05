package vane

import (
	"math"
	"net"
	"time"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/pkg/dnsutil"
	"github.com/coredns/coredns/middleware/pkg/edns"
	"github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"
	"github.com/itomsawyer/llog"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type Vane struct {
	VaneConfig

	Logger        *llog.Logger
	RcodePriority *RcodePriority
	Next          middleware.Handler
}

type VaneConfig struct {
	UpstreamTimeout time.Duration

	Debug           bool
	KeepCNAMEChain  bool
	KeepUpstreamECS bool
	AnswerShortly   bool
	ForceNoTrunc    bool
	MaxKeepA        int

	LogConfig *llog.Config
}

func NewVane() *Vane {
	return &Vane{
		RcodePriority: NewRcodePriority(),
	}
}

func (v *Vane) Init() error {
	l, err := engine.CreateLogger(v.LogConfig)
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
		cip net.IP
		i   int
		//clientSetID  int
		ok              bool
		replyMsg        *dns.Msg
		replys          []*dns.Msg
		clientSets      []engine.ClientSet
		retcode, bestrc int
		view            engine.View
		err             error
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
	v.Logger.Debug("query domain %s", q.Name)

	// Try get clientset_id from previous vane_engine middleware which has done this job.
	// In case vane_engine doesn't do its duty, try here then
	value = ctx.Value("clientsets")
	cip = state.GetRemoteAddr()
	v.Logger.Debug("get remote client addr %v", cip)
	clientSets, ok = value.([]engine.ClientSet)
	if !ok {
		if cip == nil {
			clientSets = []engine.ClientSet{engine.DefaultClientSet}
		} else {
			clientSets = e.GetClientSets(cip)
		}

		v.Logger.Debug("fetch clientsets from request: %v", clientSets)
	} else {
		v.Logger.Debug("fetch clientsets from vane_engine: %v", clientSets)
	}

	if edns.HasClientSubnet(r) && !v.KeepUpstreamECS {
		req := r.Copy()
		edns.RemoveClientSubnetIfExist(req)
		state = request.Request{W: w, Req: req}
	}

	if len(clientSets) == 0 {
		v.Logger.Warn("clientSets not found client: %v domain: %s", cip, q.Name)
		return dns.RcodeServerFailure, errUnexpectedLogic
	}

	//v.Logger.Debug("use clientsets: %v", clientSets)
	//clientSetID = clientSets[0].ID

	// Get domainpool_id , if not found e.GetDomainPoolID return engine.DefaultDomainPoolID
	domain, err := e.GetDomain(q.Name)
	if err != nil {
		domain = engine.DefaultDomainPool
	}
	origDomain := domain

try_again:
	//v.Logger.Debug("use clientset_id : %d", clientSetID)
	v.Logger.Debug("use domain pool id: %d", domain.DmPoolID)

	found_policy := make([]engine.Policy, 0, 2)
	views := make([]engine.View, 0, 2)
	for _, cs := range clientSets {
		view, err = e.GetView(cs.ID, domain.DmPoolID)
		if err != nil {
			continue
		}

		views = append(views, view)

		if view.Upstream == nil {
			continue
		}

		//Policy is a method class to choose upstreamhost (ldns) from upstream (policy)
		policy := view.Upstream.GetPolicy()
		if policy == nil {
			continue
		}

		//clientSetID = cs.ID
		found_policy = append(found_policy, policy)
	}

	// INSIST: len(found_policy) <= len(views) , we just care about found policy here
	if len(found_policy) == 0 {
		if _, domain, ok = rollback(engine.DefaultClientSetID, domain); ok {
			goto try_again
		}

		v.Logger.Warn("upstream policy not found client: %v domain: %s", cip, q.Name)
		return dns.RcodeServerFailure, errUnexpectedLogic
	}

	ex := NewExchangeHelper(nil)
	ex.Timeout = v.UpstreamTimeout
	bestrc = dns.RcodeServerFailure
	replyMsg = nil
	replys = nil

	for _, policy := range found_policy {
		for i = 0; i < defaultMaxAttampts; i++ {
			v.Logger.Debug("try select upstream hosts")
			// for each time policy choose the next prior upstreamhosts group
			uphosts := policy.Select()
			if len(uphosts) == 0 {
				// There no upstream host can be found , check out what we got.
				break
			}

			for _, uh := range uphosts {
				v.Logger.Debug("upstream host found %v", uh.Name)
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

			if v.RcodePriority.PriorTo(retcode, bestrc) {
				v.Logger.Debug("found better retcode: %d", retcode)
				bestrc = retcode
				replyMsg = replys[0]
			}

			if retcode != dns.RcodeSuccess {
				continue
			}

			// No need to filter record with type is not A, Get a proper one to return
			if q.Qtype != dns.TypeA {
				if bestrc == dns.RcodeSuccess && replyMsg != nil {
					if !v.KeepCNAMEChain {
						replyMsg.Answer = dnsutil.RemoveCNAME(replyMsg.Answer)
						replyMsg.Ns = []dns.RR{}
						replyMsg.Extra = []dns.RR{}
					}

					w.WriteMsg(replyMsg)
					return 0, nil
				}

				continue
			}

			// better is the result set of all A that pass the filter with Route
			better := make([]dns.RR, 0, 4)
			rrset := rrSet{}
			cnameSet := rrSet{}
			for _, reply := range replys {
				rrlist := reply.Answer
				for _, rr := range rrlist {
					switch r := rr.(type) {
					case (*dns.A):
						key := r.A.String()
						v.Logger.Debug("Add A %s into rrset", key)
						rrset.Add(key, r)
					case (*dns.CNAME):
						key := r.Hdr.Name
						v.Logger.Debug("Add CNAME %s -> %s into cnameset", key, r.Target)
						cnameSet.Add(key, r)
					default:
					}
				}
			}

			bestRoutePrio := math.MaxInt32
			for _, tmpView := range views {
				v.Logger.Debug("current routeset %d %s", tmpView.RouteSetID, tmpView.RouteSetName)
				for _, rr := range rrset {
					var routes engine.RouteSlice
					a := rr.(*dns.A)
					netlinks := e.GetNetLinks(a.A)
					v.Logger.Debug("ip %s found netLink: %v", a.A, netlinks)

					for _, nl := range netlinks {
						routes = e.GetRoute(tmpView.RouteSetID, domain.DmPoolID, nl.ID)
						if len(routes) > 0 {
							break
						}
					}

					if len(routes) == 0 {
						v.Logger.Debug("no route found")
						continue
					}

					route := routes[0]
					v.Logger.Debug("route found to outlink %s", route.OutLink.Addr)

					if origDomain.Monitor {
						v.Logger.Debug("domain %s with dmpool %d need to use dynamic monitor", q.Name, domain.DmPoolID)
						status, ok := e.GetLink(a.A.String(), routes[0].OutLink.Addr)
						if ok && status.Status == engine.LinkStatusDown {
							v.Logger.Debug("%s dynamic monitor status down", a)
							continue
						}
					}

					if route.Priority > bestRoutePrio {
						continue
					}

					if route.Priority < bestRoutePrio {
						bestRoutePrio = route.Priority
						better = better[:0]
						v.Logger.Debug("better route found with priority %d", route.Priority)
					}

					v.Logger.Debug("ip addr %s has been accepted for %s", a.A, route.OutLink.Addr)
					a.Hdr.Name = q.Name
					better = append(better, a)
				}

				if len(better) > 0 {
					//find at least one good  answer
					v.Logger.Debug("clientset_id %d domainpool_id %d found routes", tmpView.ClientSetID, domain.DmPoolID)
					break
				}
			}

			if len(better) == 0 {
				// nothing found, try next policy
				continue
			}

			// we got answer, return
			if v.MaxKeepA > 0 && len(better) > v.MaxKeepA {
				better = better[:v.MaxKeepA]
			}

			if cnameSlice := cnameSet.ToSlice(); len(cnameSlice) != 0 && v.KeepCNAMEChain {
				replyMsg.Answer = append(cnameSlice, better...)
			} else {
				replyMsg.Answer = better
			}

			if v.AnswerShortly {
				replyMsg.Ns = []dns.RR{}
				replyMsg.Extra = []dns.RR{}
			}

			if v.ForceNoTrunc {
				replyMsg.Truncated = false
			}

			if v.Debug {
				v.Logger.Debug("Write anwser to client: \n%v", replyMsg)
			}
			w.WriteMsg(replyMsg)
			return 0, nil
		}

		if i == defaultMaxAttampts {
			v.Logger.Warn("[%04d] policy upstream max attampts reached for client: %v domain: %s", warnPolicyMaxAttempts, cip, q.Name)
		}
	}

	// try again the whole precedure with domainPoolID equals to
	// default domainPoolID which is 1. Or the domainPoolID is already the default, Lookup failed

	if _, domain, ok = rollback(engine.DefaultClientSetID, domain); ok {
		goto try_again
	}

	// WARNING MUST be sent when we get here
	// 1. replyMsg is not nil , bestrc is dns.RcodeSuccess: we discard all of ip addr.
	// 2. domain pool has falled back to default once and still got no good answer
	if replyMsg != nil {
		if bestrc != dns.RcodeSuccess {
			//return the best effort answer we get
			w.WriteMsg(replyMsg)
			return 0, nil
		}

		// bestrc == dns.RcodeSuccess indicate that we got NOERROR response
		// but no A record we want to pick, just keep slient to the client
		v.Logger.Warn("[%04d] cannot find available answer for client: %v domain: %s", warnNoGoodAnswer, cip, q.Name)
		replyMsg.Answer = nil
		replyMsg.Ns = nil
		w.WriteMsg(replyMsg)
		return 0, nil
	}

	//we tried our best but still got nothing
	v.Logger.Warn("[%04d] cannot resolve any answer from upstream for client: %v domain: %s", warnNoneReplies, cip, q.Name)
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
