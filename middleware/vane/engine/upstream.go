package engine

import (
	"sort"

	"github.com/miekg/coredns/middleware/proxy"
)

type Upstreamer interface {
	Select() []*proxy.UpstreamHost
}

type Policy interface {
	Select(pool HostPool) []*proxy.UpstreamHost
}

type Upstream struct {
	Name   string
	Hosts  HostPool
	Policy Policy
	sorted bool
}

func NewUpstream(name string) *Upstream {
	return &Upstream{
		Name:   name,
		Hosts:  HostPool{},
		Policy: nil, //TODO defaultPolicy
	}
}

func (p *Upstream) AddHost(uh *proxy.UpstreamHost, priority int) {
	p.Hosts = p.Hosts.Add(uh, priority)
	p.sorted = false
}

func (p *Upstream) SetPolicy(policy Policy) {
	p.Policy = policy
}

func (p *Upstream) Sort() {
	sort.Sort(p.Hosts)
	p.sorted = true
}

type HostPool []HostPoolEle

func (hp HostPool) Len() int {
	return len(hp)
}

func (hp HostPool) Less(i, j int) bool {
	return hp[i].Priority < hp[j].Priority
}

func (hp HostPool) Swap(i, j int) {
	hp[i], hp[j] = hp[j], hp[i]
}

func (hp HostPool) Add(uh *proxy.UpstreamHost, priority int) HostPool {
	if hp == nil {
		hp = make([]HostPoolEle, 0, 1)
	}

	hp = append(hp, HostPoolEle{Priority: priority, Host: uh})
	return hp
}

type HostPoolEle struct {
	Priority int
	Host     *proxy.UpstreamHost
}
