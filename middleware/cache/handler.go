package cache

import (
	"math"
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

	now := c.now().UTC()
	i, ttl := c.get(now, qname, qtype, do, tag)
	if i != nil && ttl > 0 {
		resp := i.toMsg(r, uint32(c.pminttl.Seconds()))
		state.SizeAndDo(resp)
		resp, _ = state.Scrub(resp)
		w.WriteMsg(resp)

		if c.prefetch > 0 {
			i.Freq.Update(c.duration, now)

			threshold := int(math.Ceil(float64(c.percentage) / 100 * float64(i.origTTL)))

			if i.Freq.Hits() >= c.prefetch && ttl <= threshold {
				go func() {
					cachePrefetches.Inc()
					// When prefetching we loose the item i, and with it the frequency
					// that we've gathered sofar. See we copy the frequencies info back
					// into the new item that was stored in the cache.
					prr := &ResponseWriter{ResponseWriter: w, Cache: c, prefetch: true, Tag: tag}
					middleware.NextOrFailure(c.Name(), c.Next, ctx, prr, r)

					if i1 := c.exists(qname, qtype, do, tag); i1 != nil {
						i1.Freq.Reset(now, i.Freq.Hits())
					}
				}()
			}
		}

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

func (c *Cache) get(now time.Time, qname string, qtype uint16, do bool, keySuffix string) (*item, int) {
	k := rawKey(qname, qtype, do, keySuffix)

	if i, ok := c.ncache.Get(k); ok {
		cacheHits.WithLabelValues(Denial).Inc()
		return i.(*item), i.(*item).ttl(now)
	}

	if i, ok := c.pcache.Get(k); ok {
		cacheHits.WithLabelValues(Success).Inc()
		return i.(*item), i.(*item).ttl(now)
	}
	cacheMisses.Inc()
	return nil, 0
}

func (c *Cache) exists(qname string, qtype uint16, do bool, keySuffix string) *item {
	k := rawKey(qname, qtype, do, keySuffix)
	if i, ok := c.ncache.Get(k); ok {
		return i.(*item)
	}
	if i, ok := c.pcache.Get(k); ok {
		return i.(*item)
	}
	return nil
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

	cachePrefetches = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: middleware.Namespace,
		Subsystem: "cache",
		Name:      "prefetch_total",
		Help:      "The number of time the cache has prefetched a cached item.",
	})
)

const subsystem = "cache"

func init() {
	prometheus.MustRegister(cacheSize)
	prometheus.MustRegister(cacheCapacity)
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
}
