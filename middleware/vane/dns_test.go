package vane

import (
	"testing"

	"github.com/miekg/dns"
)

var (
	msgs = map[string]*dns.Msg{}
)

func init() {
	var (
		r, answer *dns.Msg
		txt       dns.RR
	)

	r = new(dns.Msg)
	r.SetQuestion("test.com", dns.TypeTXT)

	answer = new(dns.Msg)
	answer.SetRcode(r, dns.RcodeNameError)
	txt, _ = dns.NewRR(". IN 0 TXT " + "NXDOMAIN")
	answer.Answer = append(answer.Answer, txt)
	msgs["NXDOMAIN"] = answer

	answer = new(dns.Msg)
	answer.SetRcode(r, dns.RcodeServerFailure)
	txt, _ = dns.NewRR(". IN 0 TXT " + "SERVFAIL")
	answer.Answer = append(answer.Answer, txt)
	msgs["SERVFAIL"] = answer

	answer = new(dns.Msg)
	answer.SetRcode(r, dns.RcodeSuccess)
	txt, _ = dns.NewRR(". IN 0 TXT " + "NOERROR")
	answer.Answer = append(answer.Answer, txt)
	msgs["NOERROR"] = answer
}

func TestMsgSlice(t *testing.T) {
	var (
		ms, best MsgSlice
		rcode    int
	)

	ms = MsgSlice{}
	best, rcode = ms.Best()
	if rcode != dns.RcodeServerFailure {
		t.Errorf("expect best rcode %d got %d\n", dns.RcodeServerFailure, rcode)
	}

	if len(best) != 0 {
		t.Errorf("expect num of best msg is 0, got %d\n", len(best))
	}

	ms.Append(msgs["NXDOMAIN"])
	best, rcode = ms.Best()
	if rcode != dns.RcodeNameError {
		t.Errorf("expect best rcode %d got %d\n", dns.RcodeNameError, rcode)
	}

	if len(best) != 1 {
		t.Errorf("expect num of best msg is 1, got %d\n", len(best))
	}

	ms.Append(msgs["NOERROR"])
	ms.Append(msgs["NOERROR"])

	best, rcode = ms.Best()
	t.Log(best)

	if rcode != dns.RcodeSuccess {
		t.Errorf("expect best rcode %d got %d\n", dns.RcodeSuccess, rcode)
	}

	if len(best) != 2 {
		t.Errorf("expect num of best msg is 2, got %d\n", len(best))
	}
}

func TestRcodePriority(t *testing.T) {
	rp := NewRcodePriority()
	test := []struct {
		R1     int
		R2     int
		Expect bool
	}{
		{
			dns.RcodeSuccess,
			dns.RcodeSuccess,
			false,
		},
		{
			dns.RcodeSuccess,
			dns.RcodeServerFailure,
			true,
		},
		{
			dns.RcodeSuccess,
			dns.RcodeNotImplemented,
			true,
		},
		{
			dns.RcodeNameError,
			dns.RcodeServerFailure,
			true,
		},
		{
			dns.RcodeServerFailure,
			dns.RcodeNameError,
			false,
		},
		{
			dns.RcodeNXRrset,
			dns.RcodeSuccess,
			false,
		},
		{
			dns.RcodeSuccess,
			dns.RcodeNXRrset,
			true,
		},
		{
			dns.RcodeYXRrset,
			dns.RcodeNXRrset,
			false,
		},
	}

	for _, v := range test {
		if rp.PriorTo(v.R1, v.R2) != v.Expect {
			t.Errorf("expect %d prior to %d %v", v.R1, v.R2, v.Expect)
		}
	}
}
