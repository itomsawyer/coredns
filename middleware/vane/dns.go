package vane

import (
	"sort"

	"github.com/miekg/dns"
)

type RcodePriority [6]int

func (p *RcodePriority) Set(rcode int, prio int) {
	if rcode < 0 || rcode > len(p) {
		panic("RcodePriority bare code error")
	}

	p[rcode] = prio
}

// if dns retcode r1 is prior to another retcode r2.
// The priority larger, the
func (p *RcodePriority) PriorTo(r1 int, r2 int) bool {
	p1, p2 := 0, 0

	if r1 >= 0 && r1 < len(p) {
		p1 = p[r1]
	}

	if r2 >= 0 && r2 < len(p) {
		p2 = p[r2]
	}

	return p1 > p2
}

func NewRcodePriority() *RcodePriority {
	p := &RcodePriority{}
	p.Set(dns.RcodeSuccess, 10)
	p.Set(dns.RcodeNameError, 9)
	p.Set(dns.RcodeNotImplemented, 8)
	p.Set(dns.RcodeRefused, 2)
	p.Set(dns.RcodeServerFailure, 1)
	return p
}

type rrSet map[string]dns.RR

func (p rrSet) Add(key string, a dns.RR) {
	p[key] = a
}

func (p rrSet) Pack() []dns.RR {
	if len(p) == 0 {
		return nil
	}

	s := make([]dns.RR, 0, len(p))
	for _, r := range p {
		s = append(s, r)
	}

	return s
}

type MsgSlice []*dns.Msg

func (p MsgSlice) Len() int {
	return len(p)
}

func (p MsgSlice) Less(i, j int) bool {
	return p[i].Rcode < p[j].Rcode
}

func (p MsgSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *MsgSlice) Append(m *dns.Msg) {
	if p == nil {
		*p = make([]*dns.Msg, 0, 4)
	}

	*p = append(*p, m)
	return
}

func (p MsgSlice) Best() (MsgSlice, int) {
	var i int

	if len(p) == 0 {
		return nil, dns.RcodeServerFailure
	}

	if len(p) == 1 {
		return p, p[0].Rcode
	}

	sort.Sort(p)

	best := p[0].Rcode
	for i = 0; i < len(p); i++ {
		if best != p[i].Rcode {
			break
		}
	}

	return p[:i], p[0].Rcode
}
