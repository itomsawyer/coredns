package engine

import (
	"sort"
	"testing"
	"time"

	"github.com/miekg/coredns/middleware/proxy"
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

	hp.Add(c, 2)
	t.Log(hp)
	p = NewSimplePolicy(hp)
	if p == nil {
		t.Errorf("unexpected error")
		return
	}
	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())

	// change the HostPool
	hp.Add(b, 1)
	hp.Add(a, 1)
	sort.Sort(hp)
	t.Log(hp)

	p = NewSimplePolicy(hp)
	if p == nil {
		t.Errorf("unexpected error")
		return
	}

	t.Log(p.Select())
	t.Log(p.Select())
	t.Log(p.Select())
}
