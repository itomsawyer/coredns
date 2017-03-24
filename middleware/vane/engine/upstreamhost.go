package engine

import (
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware/pkg/dnsutil"
)

type UpstreamHostDownFunc func(*UpstreamHost) bool

type UpstreamHost struct {
	Conns             int64  // must be first field to be 64-bit aligned on 32-bit systems
	Name              string // IP address (and port) of this upstream host
	Fails             int32
	FailTimeout       time.Duration
	QueryTimeout      time.Duration
	Unhealthy         bool
	CheckDown         UpstreamHostDownFunc
	WithoutPathPrefix string
	Exchanger
}

// Down checks whether the upstream host is down or not.
// Down will try to use uh.CheckDown first, and will fall
// back to some default criteria if necessary.
func (uh *UpstreamHost) Down() bool {
	if uh.CheckDown == nil {
		// Default settings
		fails := atomic.LoadInt32(&uh.Fails)
		return uh.Unhealthy || fails > 0
	}
	return uh.CheckDown(uh)
}

func NewUpstreamHost(host string) *UpstreamHost {
	h, err := dnsutil.ParseHostPort(host, "53")
	if err != nil {
		return nil
	}

	uh := &UpstreamHost{
		Name:        h,
		Exchanger:   newDNSEx(),
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
