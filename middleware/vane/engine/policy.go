package engine

import (
	"github.com/miekg/coredns/middleware/proxy"
)

type Policy interface {
	Select() []*proxy.UpstreamHost
}

type PolicyBuilder func(HostPool) Policy

// XXX: Require hp to be sorted by priority
func NewSimplePolicy(hp HostPool) Policy {
	if len(hp) == 0 {
		return nil
	}

	return &SimplePolicy{
		pool: hp,
		prio: hp[0].Priority,
	}
}

type SimplePolicy struct {
	pool HostPool
	prio int
	done bool
}

func (p *SimplePolicy) Select() (uhs []*proxy.UpstreamHost) {
	var i int

	if len(p.pool) == 0 || p.done {
		return
	}

	found := false
	for i = 0; i < len(p.pool); i++ {
		e := p.pool[i]
		if e.Priority == p.prio {
			uhs = append(uhs, e.Host)
			found = true
		} else if found {
			p.prio = e.Priority
			break
		}
	}

	if i == len(p.pool) {
		p.done = true
	}

	return
}
