package vane

import (
	"github.com/miekg/dns"
)

type rrSet map[string]dns.RR

func (p rrSet) Add(a dns.RR) {
	p[a.String()] = a
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
