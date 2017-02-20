package engine

import (
	"sort"

	"github.com/miekg/coredns/middleware/proxy"
)

type Upstreamer interface {
	Select() []*proxy.UpstreamHost
	GetPolicy() Policy
}

type Upstream struct {
	Name   string
	Hosts  HostPool
	Policy PolicyBuilder
}

func NewUpstream(name string) *Upstream {
	return &Upstream{
		Name:   name,
		Hosts:  HostPool{},
		Policy: NewSimplePolicy,
	}
}

func (p *Upstream) GetPolicy() Policy {
	return p.Policy(p.Hosts)
}

func (p *Upstream) AddHost(uh *proxy.UpstreamHost, priority int) {
	p.Hosts.Add(uh, priority)
	sort.Sort(p.Hosts)
}

func (p *Upstream) SetPolicy(policy PolicyBuilder) {
	p.Policy = policy
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

func (hp *HostPool) Add(uh *proxy.UpstreamHost, priority int) {
	if hp == nil {
		*hp = make([]HostPoolEle, 0, 1)
	}

	*hp = append(*hp, HostPoolEle{Priority: priority, Host: uh})
	return
}

type HostPoolEle struct {
	Priority int
	Host     *proxy.UpstreamHost
}
