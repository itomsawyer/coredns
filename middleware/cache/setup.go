package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/hashicorp/golang-lru"
	"github.com/mholt/caddy"
)

var reuse *Cache
var lock sync.Mutex

func init() {
	caddy.RegisterPlugin("cache", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	ca, err := cacheParse(c)
	if err != nil {
		return middleware.Error("cache", err)
	}
	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		ca.Next = next
		return ca
	})

	c.OnStartup(func() error {
		lock.Lock()
		defer lock.Unlock()

		if reuse != nil {
			ca.pcache = reuse.pcache
			ca.pcache.SetLen(ca.pcap)

			ca.ncache = reuse.ncache
			ca.ncache.SetLen(ca.ncap)
		}

		reuse = ca

		// Export the capacity for the metrics. This only happens once, because this is a re-load change only.
		cacheCapacity.WithLabelValues(Success).Set(float64(ca.pcap))
		cacheCapacity.WithLabelValues(Denial).Set(float64(ca.ncap))

		return nil
	})

	return nil
}

func cacheParse(c *caddy.Controller) (*Cache, error) {

	ca := &Cache{pcap: defaultCap, ncap: defaultCap, pttl: maxTTL, nttl: maxNTTL}

	for c.Next() {
		// cache [ttl] [zones..]
		origins := make([]string, len(c.ServerBlockKeys))
		copy(origins, c.ServerBlockKeys)
		args := c.RemainingArgs()

		if len(args) > 0 {
			// first args may be just a number, then it is the ttl, if not it is a zone
			ttl, err := strconv.Atoi(args[0])
			if err == nil {
				// Reserve 0 (and smaller for future things)
				if ttl <= 0 {
					return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", ttl)
				}
				ca.pttl = time.Duration(ttl) * time.Second
				ca.nttl = time.Duration(ttl) * time.Second
				args = args[1:]
			}
			if len(args) > 0 {
				copy(origins, args)
			}
		}

		// Refinements? In an extra block.
		for c.NextBlock() {
			switch c.Val() {
			// first number is cap, second is an new ttl
			case Success:
				args := c.RemainingArgs()
				if len(args) == 0 {
					return nil, c.ArgErr()
				}
				pcap, err := strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
				ca.pcap = pcap
				if len(args) > 1 {
					pttl, err := strconv.Atoi(args[1])
					if err != nil {
						return nil, err
					}
					// Reserve 0 (and smaller for future things)
					if pttl <= 0 {
						return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", pttl)
					}
					ca.pttl = time.Duration(pttl) * time.Second
				}
			case Denial:
				args := c.RemainingArgs()
				if len(args) == 0 {
					return nil, c.ArgErr()
				}
				ncap, err := strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
				ca.ncap = ncap
				if len(args) > 1 {
					nttl, err := strconv.Atoi(args[1])
					if err != nil {
						return nil, err
					}
					// Reserve 0 (and smaller for future things)
					if nttl <= 0 {
						return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", nttl)
					}
					ca.nttl = time.Duration(nttl) * time.Second
				}
			default:
				return nil, c.ArgErr()
			}
		}

		for i := range origins {
			origins[i] = middleware.Host(origins[i]).Normalize()
		}

		var err error
		ca.Zones = origins

		ca.pcache, err = lru.New(ca.pcap)
		if err != nil {
			return nil, err
		}
		ca.ncache, err = lru.New(ca.ncap)
		if err != nil {
			return nil, err
		}

		return ca, nil
	}

	return nil, nil
}
