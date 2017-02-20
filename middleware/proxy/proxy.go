// Package proxy is middleware that proxies requests.
package proxy

import (
	"errors"
	"sync/atomic"
	"time"

	"github.com/miekg/coredns/middleware"
	"github.com/miekg/coredns/middleware/pkg/dnsutil"
	"github.com/miekg/coredns/request"

	"github.com/miekg/dns"
	ot "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

var (
	errUnreachable     = errors.New("unreachable backend")
	errInvalidProtocol = errors.New("invalid protocol")
)

// Proxy represents a middleware instance that can proxy requests to another (DNS) server.
type Proxy struct {
	Next      middleware.Handler
	Upstreams []Upstream
}

// Upstream manages a pool of proxy upstream hosts. Select should return a
// suitable upstream host, or nil if no such hosts are available.
type Upstream interface {
	// The domain name this upstream host should be routed on.
	From() string
	// Selects an upstream host to be routed to.
	Select() *UpstreamHost
	// Checks if subpdomain is not an ignored.
	IsAllowedPath(string) bool
	// Options returns the options set for this upstream
	Options() Options
}

// UpstreamHostDownFunc can be used to customize how Down behaves.
type UpstreamHostDownFunc func(*UpstreamHost) bool

// UpstreamHost represents a single proxy upstream
type UpstreamHost struct {
	Conns             int64  // must be first field to be 64-bit aligned on 32-bit systems
	Name              string // IP address (and port) of this upstream host
	Fails             int32
	FailTimeout       time.Duration
	QueryTimeout      time.Duration
	Unhealthy         bool
	CheckDown         UpstreamHostDownFunc
	WithoutPathPrefix string
	Exchanger
}

func NewUpstreamHost(host string, timeout time.Duration) *UpstreamHost {
	h, err := dnsutil.ParseHostPort(host, "53")
	if err != nil {
		return nil
	}

	uh := &UpstreamHost{
		Name:        h,
		Conns:       0,
		Fails:       0,
		FailTimeout: 10 * time.Second,
		Unhealthy:   false,
		Exchanger:   newDNSEx(h),

		CheckDown: func(uh *UpstreamHost) bool {
			if uh.Unhealthy {
				return true
			}

			fails := atomic.LoadInt32(&uh.Fails)
			if fails >= 2 {
				return true
			}
			return false
		},
		WithoutPathPrefix: "",
	}

	uh.SetTimeout(timeout)
	return uh
}

func (uh *UpstreamHost) SetQueryTimeout(timeout time.Duration) {
	uh.QueryTimeout = timeout
	uh.SetTimeout(timeout)
}

func (uh *UpstreamHost) DoExchange(state request.Request) (*dns.Msg, error) {
	atomic.AddInt64(&uh.Conns, 1)
	defer atomic.AddInt64(&uh.Conns, -1)
	return uh.Exchange(state)
}

func (uh *UpstreamHost) Fail() {
	timeout := uh.FailTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	atomic.AddInt32(&uh.Fails, 1)
	go func(host *UpstreamHost, timeout time.Duration) {
		time.Sleep(timeout)
		atomic.AddInt32(&host.Fails, -1)
	}(uh, timeout)
}

// Down checks whether the upstream host is down or not.
// Down will try to use uh.CheckDown first, and will fall
// back to some default criteria if necessary.
func (uh *UpstreamHost) Down() bool {
	if uh.CheckDown == nil {
		// Default settings
		fails := atomic.LoadInt32(&uh.Fails)
		return uh.Unhealthy || fails > 0
	}
	return uh.CheckDown(uh)
}

// tryDuration is how long to try upstream hosts; failures result in
// immediate retries until this duration ends or we get a nil host.
var tryDuration = 60 * time.Second

// ServeDNS satisfies the middleware.Handler interface.
func (p Proxy) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var (
		span, child ot.Span
	)

	span = ot.SpanFromContext(ctx)
	state := request.Request{W: w, Req: r}

	for _, upstream := range p.Upstreams {
		start := time.Now()

		// Since Select() should give us "up" hosts, keep retrying
		// hosts until timeout (or until we get a nil host).
		for time.Now().Sub(start) < tryDuration {
			host := upstream.Select()
			if host == nil {
				RequestDuration.WithLabelValues(upstream.From()).Observe(float64(time.Since(start) / time.Millisecond))
				return dns.RcodeServerFailure, errUnreachable
			}

			if span != nil {
				child = span.Tracer().StartSpan("exchange", ot.ChildOf(span.Context()))
				ctx = ot.ContextWithSpan(ctx, child)
			}

			atomic.AddInt64(&host.Conns, 1)

			reply, backendErr := host.Exchange(state)

			atomic.AddInt64(&host.Conns, -1)

			if child != nil {
				child.Finish()
			}

			if backendErr == nil {
				w.WriteMsg(reply)
				RequestDuration.WithLabelValues(upstream.From()).Observe(float64(time.Since(start) / time.Millisecond))
				return 0, nil
			}
			timeout := host.FailTimeout
			if timeout == 0 {
				timeout = 10 * time.Second
			}
			atomic.AddInt32(&host.Fails, 1)
			go func(host *UpstreamHost, timeout time.Duration) {
				time.Sleep(timeout)
				atomic.AddInt32(&host.Fails, -1)
			}(host, timeout)
		}

		RequestDuration.WithLabelValues(upstream.From()).Observe(float64(time.Since(start) / time.Millisecond))

		return dns.RcodeServerFailure, errUnreachable
	}
	return middleware.NextOrFailure(p.Name(), p.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (p Proxy) Name() string { return "proxy" }

// defaultTimeout is the default networking timeout for DNS requests.
const defaultTimeout = 5 * time.Second
