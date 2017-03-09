package proxy

import (
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware/pkg/dnsutil"
)

func NewUpstreamHost(host string) *UpstreamHost {
	h, err := dnsutil.ParseHostPort(host, "53")
	if err != nil {
		return nil
	}

	uh := &UpstreamHost{
		Name:        h,
		Conns:       0,
		Fails:       0,
		FailTimeout: 10 * time.Second,
		Unhealthy:   false,
		CheckDown:   nil,
		/*
			CheckDown: func(uh *UpstreamHost) bool {
				if uh.Unhealthy {
					return true
				}

				fails := atomic.LoadInt32(&uh.Fails)
				if fails >= 3 {
					return true
				}
				return false
			},
		*/
		WithoutPathPrefix: "",
	}

	return uh
}

func (uh *UpstreamHost) Fail() {
	timeout := uh.FailTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	atomic.AddInt32(&uh.Fails, 1)
	go func(host *UpstreamHost, timeout time.Duration) {
		time.Sleep(timeout)
		atomic.AddInt32(&host.Fails, -1)
	}(uh, timeout)
}
