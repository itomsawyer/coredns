package request

import (
	"testing"

	"github.com/coredns/coredns/middleware/test"

	"github.com/miekg/dns"
)

func TestRequestDo(t *testing.T) {
	st := testRequest()

	st.Do()
	if st.do == 0 {
		t.Fatalf("Expected st.do to be set")
	}
}

func TestRequestRemote(t *testing.T) {
	st := testRequest()
	if st.IP() != "10.240.0.1" {
		t.Fatalf("Wrong IP from request")
	}
	p := st.Port()
	if p == "" {
		t.Fatalf("Failed to get Port from request")
	}
	if p != "40212" {
		t.Fatalf("Wrong port from request")
	}
}

func TestResponse(t *testing.T) {
	rt := testResponse()
	ans := rt.GetResponseSummaryText()
	for _, rr := range rt.Req.Answer {
		t.Log("test reponse answer", rr)
	}

	if len(ans) != 2 {
		t.Errorf("ans length not expected")
		return
	}

	if ans[0] != "epl.com." {
		t.Errorf("ans record not expected")
		t.Log(ans[0])
		return
	}

	if ans[1] != "1.1.1.1" {
		t.Errorf("ans record not expected")
		t.Log(ans[1])
		return
	}
}

func BenchmarkRequestDo(b *testing.B) {
	st := testRequest()

	for i := 0; i < b.N; i++ {
		st.Do()
	}
}

func BenchmarkRequestSize(b *testing.B) {
	st := testRequest()

	for i := 0; i < b.N; i++ {
		st.Size()
	}
}

func testRequest() Request {
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)
	m.SetEdns0(4097, true)
	return Request{W: &test.ResponseWriter{}, Req: m}
}

func testResponse() Request {
	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeA)
	m.SetEdns0(4097, true)
	rr1, _ := dns.NewRR("example.com. 300 CNAME epl.com.")
	rr2, _ := dns.NewRR("example.com. 300 A 1.1.1.1")

	m.Answer = append(m.Answer, rr1)
	m.Answer = append(m.Answer, rr2)

	return Request{W: &test.ResponseWriter{}, Req: m}
}
