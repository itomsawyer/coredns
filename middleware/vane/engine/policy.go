package engine

import (
	"github.com/coredns/coredns/middleware/proxy"
)

type Policy interface {
	Select() []*proxy.UpstreamHost
}

type PolicyBuilder func(HostPool) Policy

func NewSimplePolicy(hp HostPool) Policy {
	if len(hp) == 0 {
		return nil
	}

	return &SimplePolicy{
		pool: hp,
		//prio: hp[0].Priority,
		cur: 0,
		len: len(hp),
	}
}

type SimplePolicy struct {
	pool HostPool

	cur int
	len int
}

// Require HostPool to be sorted by priority
func (p *SimplePolicy) Select() (uhs []*proxy.UpstreamHost) {
	var i int

	if p.cur < 0 || p.cur >= p.len {
		return nil
	}

	found := false
	for i = p.cur; i < p.len; i++ {
		if p.pool[i].Priority != p.pool[p.cur].Priority {
			p.cur = i
			if found {
				return
			} else {
				i--
				continue
			}
		}

		uh := p.pool[i].Host
		if uh == nil || uh.Down() {
			continue
		}

		uhs = append(uhs, uh)
		found = true
	}

	if i == p.len {
		p.cur = i
	}

	return

}
