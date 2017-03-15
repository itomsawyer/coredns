package vane

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/middleware/vane/engine"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type ExchangeHelper struct {
	Upstream *engine.Upstream
	Hosts    []*proxy.UpstreamHost
	Timeout  time.Duration
}

func NewExchangeHelper(u *engine.Upstream, hosts []*proxy.UpstreamHost) *ExchangeHelper {
	return &ExchangeHelper{
		Upstream: u,
		Hosts:    hosts,
	}
}

// replys are the dns responses from upstream dns server, max length is equal to len(h.Hosts)
// all of replys will be sorted by retcode to make sure the better reponse always comes first
// retcode NOERROR PriorTo NXDOMAIN PriorTo NOTIMPLEMENT PriorTo REFUSE PriorTo SERVERFAIL
// see msgs.Best() for details
func (h *ExchangeHelper) DoExchange(ctx context.Context, state request.Request) (replys []*dns.Msg, retcode int) {
	if h.Upstream == nil {
		return nil, dns.RcodeServerFailure
	}

	ex := h.Upstream.Exchanger()
	if ex == nil {
		return nil, dns.RcodeServerFailure
	}

	if len(h.Hosts) == 1 {
		uh := h.Hosts[0]

		atomic.AddInt64(&uh.Conns, 1)
		reply, backendErr := ex.Exchange(ctx, uh.Name, state)
		atomic.AddInt64(&uh.Conns, -1)

		//XXX gaoxiang if uh.Fail() is needed? if vane should check upstream ldns
		//    with gwCheck is working
		if backendErr != nil {
			//uh.Fail()
			return nil, dns.RcodeServerFailure
		}

		return []*dns.Msg{reply}, reply.Rcode
	}

	msgs := MsgSlice{}
	out := make(chan *dns.Msg)
	errChan := make(chan error)
	wg := new(sync.WaitGroup)
	done := make(chan struct{})

	for i := 0; i < len(h.Hosts); i++ {
		wg.Add(1)
		go func(uh *proxy.UpstreamHost) {
			atomic.AddInt64(&uh.Conns, 1)
			reply, backendErr := ex.Exchange(ctx, uh.Name, state)
			atomic.AddInt64(&uh.Conns, -1)

			if backendErr == nil {
				select {
				case out <- reply:
				case <-done:
				}
			} else {
				select {
				case errChan <- backendErr:
				case <-done:
				}
				//XXX gaoxiang if uh.Fail() is needed? if vane should check upstream ldns
				//    with gwCheck is working
				//uh.Fail()
			}

			wg.Done()
		}(h.Hosts[i])
	}

	defer func() {
		close(done)
		go func() {
			wg.Wait()
			close(out)
			close(errChan)
		}()
	}()

	if h.Timeout == 0 {
		h.Timeout = defaultTimeout
	}

	t := time.NewTimer(h.Timeout)

	for cnt := 0; cnt < len(h.Hosts); cnt++ {
		select {
		case reply := <-out:
			msgs.Append(reply)
		case <-errChan:
		case <-t.C:
			return
		}
	}

	return msgs.Best()
}
