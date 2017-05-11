package dnsutil

import "github.com/miekg/dns"

// DuplicateCNAME returns true if r already exists in records.
func DuplicateCNAME(r *dns.CNAME, records []dns.RR) bool {
	for _, rec := range records {
		if v, ok := rec.(*dns.CNAME); ok {
			if v.Target == r.Target {
				return true
			}
		}
	}
	return false
}

func RemoveCNAME(records []dns.RR) []dns.RR {
	left := make([]dns.RR, 0, 1)

	for _, rr := range records {
		if _, ok := rr.(*dns.CNAME); ok {
			continue
		}

		left = append(left, rr)
	}

	return left
}
