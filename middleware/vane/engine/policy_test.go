package engine

import (
	"sort"
	"testing"
	"time"

	"github.com/coredns/coredns/middleware/proxy"
)

func TestSimplePolicy(t *testing.T) {
	a := proxy.NewUpstreamHost("1.1.1.1", 1*time.Second)
	b := proxy.NewUpstreamHost("1.1.1.2", 1*time.Second)
	c := proxy.NewUpstreamHost("1.1.1.3", 1*time.Second)

	hp := HostPool{}

	p := NewSimplePolicy(hp)
	if p != nil {
		t.Errorf("unexpected no error")
		return
	}

	hp.Add(c, 3)
	t.Log("hostpool", hp)
	p = NewSimplePolicy(hp)
	if p == nil {
		t.Errorf("unexpected error")
		return
	}
	t.Log("policy", p)
	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())

	// change the HostPool
	b.Unhealthy = true
	hp.Add(b, 2)
	hp.Add(a, 1)
	sort.Sort(hp)
	t.Log("hostpool", hp)

	p = NewSimplePolicy(hp)
	if p == nil {
		t.Errorf("unexpected error")
		return
	}

	t.Log("policy", p)
	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())
}
