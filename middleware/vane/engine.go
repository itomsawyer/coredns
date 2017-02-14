package vane

import (
	"errors"
	"net"

	"github.com/miekg/coredns/middleware/pkg/dmtree"
	"github.com/miekg/coredns/middleware/pkg/iptree"
	"github.com/miekg/coredns/middleware/proxy"
)

type Enginer interface {
	//GetClientSetID(ip net.IP) (int, error)
	//GetDomainID(domain string) (int, error)
	//GetPolicy(client int, domain int) (policy int, router int, err error)
	//GetOutLink(router int, domain int, netlink int) error
}

type Engine struct {
	L         Loader
	Loaded    bool
	ClientSet *iptree.IPTree
	NetLink   *iptree.IPTree
	DomainSet *dmtree.DmTree
	Proxy     map[int]*proxy.Proxy
}

func NewEngine(l Loader) *Engine {
	return &Engine{
		L:      l,
		Loaded: false,
	}
}

func (e *Engine) GetClientSetID(ip net.IP) (int, error) {
	if ip == nil {
		return 0, errors.New("client ip is nil")
	}

	id, found, err := e.ClientSet.Get(ip)
	if !found {
		return 0, errors.New("ClientSet not found")
	}
	return id, err
}

func (e *Engine) GetDomainID(domain string) (int, error) {
	v, ok := e.DomainSet.Find(domain)
	if !ok && v == nil {
		return 0, errors.New("domain pool id not found")
	}

	return v.(int), nil
}
