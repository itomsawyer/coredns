package engine

import (
	"testing"

	"github.com/coredns/coredns/middleware/proxy"
)

func TestAddUpstream(t *testing.T) {
	upstream := NewUpstream("test")

	a := proxy.NewUpstreamHost("1.1.1.1")
	b := proxy.NewUpstreamHost("1.1.1.2")
	c := proxy.NewUpstreamHost("1.1.1.3")

	upstream.AddHost(b, 2)
	upstream.AddHost(c, 3)
	upstream.AddHost(a, 1)

	t.Log(upstream)

	p := upstream.GetPolicy()
	if p == nil {
		t.Errorf("unexpected nil policy")
	}

	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())
}
