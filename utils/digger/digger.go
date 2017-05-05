package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

var domains []string

func main() {
	var o *os.File
	var err error

	conFile := flag.String("conf", "domains.conf", "domains to lookup (JSON format)")
	host := flag.String("host", "223.5.5.5:53", "dns server host")
	output := flag.String("output", "stdout", "output file")
	hasCname := flag.Bool("hascname", true, "if true, dns response must has CNAME recored to be considered valid response")
	hasA := flag.Bool("hasa", true, "if true, dns response must has A recored to be considered valid response")

	flag.Parse()

	conf, err := ioutil.ReadFile(*conFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = json.Unmarshal(conf, &domains); err != nil {
		fmt.Println(err)
		return
	}

	switch *output {
	case "stdout":
		o = os.Stdout
	case "stderr":
		o = os.Stderr
	default:
		o, err = os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer o.Close()
	}

	v := Validator{}
	if *hasA {
		v.AddFilter(NewHasAFilter(true))
	}
	if *hasCname {
		v.AddFilter(NewHasCNAMEFilter(true))
	}

	c := new(dns.Client)
	for _, dm := range domains {
		msg := newQuestion(dm)
		if msg == nil {
			continue
		}

		in, rtt, err := c.Exchange(msg, *host)
		if err != nil {
			ret := fmt.Sprintf("[ErrExchange:%s]", err.Error())
			fmt.Fprintf(o, "%s %s %s %d %s\n", dm, "-", rtt, rtt, ret)
			continue
		}

		ret := "[OK]"
		if err := v.Check(MsgContent{Msg: in, RTT: rtt}); err != nil {
			ret = fmt.Sprintf("[%s]", err.Error())
		}

		records := ResponseSummaryText(in)
		fmt.Fprintf(o, "%s %s %s %d %s\n", dm, records.String(), rtt, rtt, ret)
	}
}

func newQuestion(dm string) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(dm), dns.TypeA)
	m.SetEdns0(dns.DefaultMsgSize, false)
	return m
}

type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) String() string {
	if len(p) == 0 {
		return "-"
	}

	return strings.Join(p, "|")
}

type MsgContent struct {
	Msg *dns.Msg
	RTT time.Duration
}

type Validator struct {
	Filters []ValidateFunc
}

func (v *Validator) AddFilter(f ValidateFunc) {
	if v.Filters == nil {
		v.Filters = make([]ValidateFunc, 0, 1)
	}

	v.Filters = append(v.Filters, f)
}

func (v *Validator) Check(m MsgContent) error {
	for _, f := range v.Filters {
		if err := f(m); err != nil {
			return err
		}
	}

	return nil
}

type ValidateFunc func(MsgContent) error

func NewHasAFilter(expect bool) ValidateFunc {
	return func(m MsgContent) error {
		var has bool
		for _, rr := range m.Msg.Answer {
			if _, ok := rr.(*dns.A); ok {
				has = true
				break
			}
		}

		if has != expect {
			return fmt.Errorf("ErrHasA:%v", has)
		}
		return nil
	}
}

func NewHasCNAMEFilter(expect bool) ValidateFunc {
	return func(m MsgContent) error {
		var has bool
		for _, rr := range m.Msg.Answer {
			if _, ok := rr.(*dns.CNAME); ok {
				has = true
				break
			}
		}

		if has != expect {
			return fmt.Errorf("ErrHasCNAME:%v", has)
		}
		return nil
	}
}

func ResponseSummaryText(resp *dns.Msg) StringSlice {
	answer := StringSlice{}
	for _, rr := range resp.Answer {
		if a, ok := rr.(*dns.A); ok {
			answer = append(answer, a.A.String())
		}
		if cname, ok := rr.(*dns.CNAME); ok {
			answer = append(answer, cname.Target)
		}
	}

	sort.Sort(answer)
	return answer
}
