package cache

import (
	"strconv"
	"time"

	"github.com/coredns/coredns/middleware"
	vane "github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

// ServeDNS implements the middleware.Handler interface.
func (c *Cache) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	v := ctx.Value("vane_engine")
	tag := "0"
	if engine, ok := v.(*vane.Engine); ok && engine != nil {
		cip := state.GetRemoteAddr()
		if cip != nil {
			clientSets := engine.GetClientSets(cip)
			if len(clientSets) > 0 {
				cid := clientSets[0].ID
				tag = strconv.Itoa(cid)
				ctx = context.WithValue(ctx, "clientsets", clientSets)
			}
		}
	}

	qname := state.Name()
	qtype := state.QType()
	zone := middleware.Zones(c.Zones).Matches(qname)
	if zone == "" {
		return c.Next.ServeDNS(ctx, w, r)
	}

	do := state.Do() // TODO(): might need more from OPT record? Like the actual bufsize?

	if i, ok, expired := c.get(qname, qtype, do, tag); ok && !expired {
		resp := i.toMsg(r)
		state.SizeAndDo(resp)
		resp, _ = state.Scrub(resp)
		w.WriteMsg(resp)

		return dns.RcodeSuccess, nil
	}

	crr := &ResponseWriter{
		ResponseWriter: w,
		Cache:          c,
		Tag:            tag,
	}
	return middleware.NextOrFailure(c.Name(), c.Next, ctx, crr, r)
}

// Name implements the Handler interface.
func (c *Cache) Name() string { return "cache" }

func (c *Cache) get(qname string, qtype uint16, do bool, keySuffix string) (*item, bool, bool) {
	k := rawKey(qname, qtype, do, keySuffix)

	if i, ok := c.ncache.Get(k); ok {
		cacheHits.WithLabelValues(Denial).Inc()
		return i.(*item), ok, i.(*item).expired(time.Now())
	}

	if i, ok := c.pcache.Get(k); ok {
		cacheHits.WithLabelValues(Success).Inc()
		return i.(*item), ok, i.(*item).expired(time.Now())
	}
	cacheMisses.Inc()
	return nil, false, false
}

var (
	cacheSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: middleware.Namespace,
		Subsystem: subsystem,
		Name:      "size",
		Help:      "The number of elements in the cache.",
	}, []string{"type"})

	cacheCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: middleware.Namespace,
		Subsystem: subsystem,
		Name:      "capacity",
		Help:      "The cache's capacity.",
	}, []string{"type"})

	cacheHits = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: middleware.Namespace,
		Subsystem: subsystem,
		Name:      "hits_total",
		Help:      "The count of cache hits.",
	}, []string{"type"})

	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: middleware.Namespace,
		Subsystem: subsystem,
		Name:      "misses_total",
		Help:      "The count of cache misses.",
	})
)

const subsystem = "cache"

func init() {
	prometheus.MustRegister(cacheSize)
	prometheus.MustRegister(cacheCapacity)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
}
