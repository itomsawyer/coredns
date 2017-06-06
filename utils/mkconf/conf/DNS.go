package conf

import (
	"bytes"
	"sort"
)

type DNS struct {
	Name        string
	Host        []string
	CheckDomain []string
	ExForwarder []string
}

func (p *DNS) footprint() string {
	buf := new(bytes.Buffer)
	tmp := make([]string, len(p.Host), len(p.Host))
	copy(tmp, p.Host)

	sort.Slice(tmp, func(i, j int) bool { return tmp[i] < tmp[j] })
	for _, m := range tmp {
		buf.WriteString(m + ";")
	}

	return buf.String()
}

func (p *DNS) AddHost(dns string, checkdm string) {
	if p.Host == nil {
		p.Host = make([]string, 0, 1)
	}
	p.Host = append(p.Host, dns)

	if p.CheckDomain == nil {
		p.CheckDomain = make([]string, 0, 1)
	}

	p.CheckDomain = append(p.CheckDomain, checkdm)

	p.Name = p.footprint()
}

func (p *DNS) AddExForwarder(dns string) {
	if p.ExForwarder == nil {
		p.ExForwarder = make([]string, 0, 1)
	}

	p.ExForwarder = append(p.ExForwarder, dns)
}
