package engine

import (
	"errors"
	"net"
	"time"

	"github.com/miekg/coredns/middleware/pkg/dmtree"
	"github.com/miekg/coredns/middleware/pkg/dnsutil"
	"github.com/miekg/coredns/middleware/pkg/iptree"
	"github.com/miekg/coredns/middleware/proxy"
)

var (
	DefaultClientSetID  = 1
	DefaultDomainPoolID = 1
	DefaultNetLinkID    = 1
)

/*
type Enginer interface {
	GetClientSetID(ip net.IP) (int, error)
	GetDomainPoolID(domain string) (int, error)
	GetNetLinkID(ip net.IP) (int, error)
	GetRouteSetID(clientset, dmpool int) (int, error)
	GetNetLinkSetID(netlink, dmpool int) (int, error)
	FilterTarget(netlinkset int)
	SelectLDNS(clientset, dmpool int) ([]*proxy.UpstreamHost, error)
}
*/

type Engine struct {
	ClientSet *iptree.IPTree
	NetLink   *iptree.IPTree
	Domain    *dmtree.DmTree

	UpstreamHosts map[string]*proxy.UpstreamHost
	Upstream      map[int]*Upstream

	SrcView // <dmpoolid, clientsetid> => {Route, Upstream}
	DstView // <dmpoolid, netlintid> => netlinksetid

	RouteMap // <routeid, netlinksetid> => Route
}

func (e *Engine) GetClientSetID(ip net.IP) (clientset_id int) {
	cs, err := e.GetClientSet(ip)
	if err != nil {
		return DefaultClientSetID
	}

	if cs.ID > 0 {
		return cs.ID
	}

	return DefaultClientSetID
}

func (e *Engine) GetClientSet(ip net.IP) (ClientSet, error) {
	if ip == nil {
		return ClientSet{}, errors.New("client ip is nil")
	}

	v, found, _ := e.ClientSet.GetRaw(ip)
	if !found {
		return ClientSet{}, errors.New("ClientSet not found")
	}

	cs, ok := v.(ClientSet)
	if !ok {
		return ClientSet{}, errors.New("Internal Error: ClientSet type assert fail")
	}

	return cs, nil
}

func (e *Engine) GetNetLinkID(ip net.IP) int {
	nl, err := e.GetNetLink(ip)
	if err != nil {
		return DefaultNetLinkID
	}

	if nl.ID <= 0 {
		return DefaultNetLinkID
	}

	return nl.ID
}

func (e *Engine) GetNetLink(ip net.IP) (NetLink, error) {
	if ip == nil {
		return NetLink{}, errors.New("target ip is nil")
	}

	v, found, _ := e.NetLink.GetRaw(ip)
	if !found {
		return NetLink{}, errors.New("NetLink not found")
	}

	nl, ok := v.(NetLink)
	if !ok {
		return NetLink{}, errors.New("Internal Error: NetLink type assert fail")
	}

	return nl, nil
}

func (e *Engine) GetDomainPoolID(domain string) int {
	dm, err := e.GetDomain(domain)
	if err != nil {
		return DefaultDomainPoolID
	}

	if dm.DmPoolID <= 0 {
		return DefaultDomainPoolID
	}

	return dm.DmPoolID
}

func (e *Engine) GetDomain(domain string) (Domain, error) {
	v, ok := e.Domain.Find(domain)
	if !ok && v == nil {
		return Domain{}, errors.New("domain pool id not found")
	}

	d, ok := v.(Domain)
	if !ok {
		return Domain{}, errors.New("type of domain is invalid")
	}

	return d, nil
}

func (e *Engine) GetRoute(routeset_id int, domainpool_id int, netlink_id int) RouteSlice {
	dl, err := e.GetDomainLink(domainpool_id, netlink_id)
	if err != nil {
		return nil
	}

	rk := RouteKey{
		RouteSetID:   routeset_id,
		NetLinkSetID: dl.NetLinkSetID,
	}
	if v, ok := e.RouteMap[rk]; ok {
		return v
	}

	return nil
}

func (e *Engine) GetView(clientset_id int, domainpool_id int) (View, error) {
	vk := ViewKey{
		ClientSetID:  clientset_id,
		DomainPoolID: domainpool_id,
	}

	if v, ok := e.SrcView[vk]; ok {
		return v, nil
	}

	return View{}, errors.New("View not found")
}

func (e *Engine) GetDomainLink(domainpool_id, netlink_id int) (DomainLink, error) {
	dlk := DomainLinkKey{
		DomainPoolID: domainpool_id,
		NetLinkID:    netlink_id,
	}

	if dl, ok := e.DstView[dlk]; ok {
		return dl, nil
	}

	return DomainLink{}, errors.New("DomainLink not found")
}

func (e *Engine) AddClient(ipnet *net.IPNet, id int, name string) error {
	cs := ClientSet{
		ID:    id,
		Name:  name,
		IPNet: ipnet,
	}
	if e.ClientSet == nil {
		e.ClientSet = iptree.New()
	}

	return e.ClientSet.AddRaw(cs.IPNet, cs)
}

func (e *Engine) AddNetLink(ipnet *net.IPNet, id int, isp string, region string) error {
	nl := NetLink{
		ID:     id,
		Isp:    isp,
		Region: region,
		IPNet:  ipnet,
	}

	if e.NetLink == nil {
		e.NetLink = iptree.New()
	}

	return e.NetLink.AddRaw(nl.IPNet, nl)
}

func (e *Engine) AddDomain(id int, domain string, dmpool_id int, dmpool_name string) error {
	d := Domain{
		ID:       id,
		Domain:   domain,
		DmPoolID: dmpool_id,
		DmPool:   dmpool_name,
	}

	if e.Domain == nil {
		e.Domain = new(dmtree.DmTree)
	}

	return e.Domain.Insert(d.Domain, d)
}

func (e *Engine) AddUpstreamHost(host string, timeout time.Duration, unhealthy bool) (*proxy.UpstreamHost, error) {
	h, err := dnsutil.ParseHostPort(host, "53")
	if err != nil {
		return nil, err
	}

	if e.UpstreamHosts == nil {
		e.UpstreamHosts = make(map[string]*proxy.UpstreamHost, 8)
	}

	if uh, ok := e.UpstreamHosts[h]; ok {
		uh.SetTimeout(timeout)
		uh.Unhealthy = unhealthy
		return uh, nil
	}

	uh := proxy.NewUpstreamHost(h, timeout)
	if uh == nil {
		return nil, errors.New("upsteam host (ldns) address format error")
	}
	uh.Unhealthy = unhealthy
	e.UpstreamHosts[uh.Name] = uh

	return uh, nil
}

func (e *Engine) AddUpstream(policy int, name string) (*Upstream, error) {
	if e.Upstream == nil {
		e.Upstream = make(map[int]*Upstream, 8)
	}
	u := NewUpstream(name)
	e.Upstream[policy] = u
	return u, nil
}

func (e *Engine) GetUpstreamByID(policy int) (*Upstream, error) {
	if e.Upstream == nil {
		return nil, errors.New("upstream (policy) not found")
	}

	u, ok := e.Upstream[policy]
	if !ok || u == nil {
		return nil, errors.New("upstream (policy) not found")
	}

	return u, nil
}

func (e *Engine) AttachUpstreamHost(policy int, host *proxy.UpstreamHost, priority int) error {
	upstream, err := e.GetUpstreamByID(policy)
	if err != nil {
		return err
	}

	upstream.AddHost(host, priority)
	return nil
}

func (e *Engine) AddDomainLink(dl DomainLink) {
	if e.DstView == nil {
		e.DstView = DstView{}
	}

	e.DstView.AddDomainLink(dl)
}

func (e *Engine) AddView(view View) {
	if e.SrcView == nil {
		e.SrcView = SrcView{}
	}

	e.SrcView.AddView(view)
}

func (e *Engine) AddRoute(route Route) {
	if e.RouteMap == nil {
		e.RouteMap = RouteMap{}
	}

	e.RouteMap.AddRoute(route)
}
