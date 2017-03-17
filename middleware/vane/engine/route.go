package engine

import (
	"errors"
	"net"
	"sort"
)

type RouteMap map[RouteKey]RouteSlice

func (m RouteMap) AddRoute(r Route) {
	v, ok := m[r.RouteKey]
	if !ok {
		m[r.RouteKey] = RouteSlice{r}
		return
	}

	v = append(v, r)
	sort.Sort(v)
	m[r.RouteKey] = v
}

type Route struct {
	RouteKey
	OutLink  OutLink
	Priority int
	Score    int
}

func NewRoute(routeset_id int, netlinkset_id int, o OutLink) Route {
	return Route{
		RouteKey: RouteKey{RouteSetID: routeset_id, NetLinkSetID: netlinkset_id},
		OutLink:  o,
	}
}

type RouteKey struct {
	RouteSetID   int
	NetLinkSetID int
}

type OutLink struct {
	Name string
	Addr net.IP
}

func NewOutLink(name, ipaddr string) (o OutLink, err error) {
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		err = errors.New("Ip address format error")
		return
	}

	return OutLink{Name: name, Addr: ip.To4()}, nil
}

type RouteSlice []Route

func (r RouteSlice) Len() int { return len(r) }
func (r RouteSlice) Less(i, j int) bool {
	if r[i].Priority != r[j].Priority {
		return r[i].Priority < r[j].Priority
	}
	return r[i].Score > r[j].Score
}

func (r RouteSlice) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

func (r RouteSlice) Iterator() RouteIter {
	return RouteIter{
		route: r,
		cur:   0,
		len:   len(r),
	}
}

type RouteIter struct {
	route RouteSlice
	cur   int
	len   int
}

// Require RouteSlice to be sorted which is ensured by RouteMap.Add
func (ri *RouteIter) Next() []Route {
	if ri.cur < 0 || ri.cur >= ri.len {
		return nil
	}

	end := 0
	for end = ri.cur; end < ri.len; end++ {
		if ri.route[end].Priority != ri.route[ri.cur].Priority {
			break
		}
	}

	route := ri.route[ri.cur:end]
	ri.cur = end
	return route
}
