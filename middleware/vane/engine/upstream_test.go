package engine

import (
	"testing"
	"time"

	"github.com/miekg/coredns/middleware/proxy"
)

func TestAddUpstream(t *testing.T) {
	upstream := NewUpstream("test")

	a := proxy.NewUpstreamHost("1.1.1.1", 1*time.Second)
	b := proxy.NewUpstreamHost("1.1.1.2", 1*time.Second)
	c := proxy.NewUpstreamHost("1.1.1.3", 1*time.Second)

	upstream.AddHost(b, 2)
	upstream.AddHost(c, 3)
	upstream.AddHost(a, 1)

	t.Log(upstream)
}
