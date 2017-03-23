package engine

import (
	"time"
)

type Dst2LNK struct {
	DstIP   string `json:"dst_ip"`
	OutLink string `json:"out_link"`
}

const (
	LinkStatusUnknown = -1
	LinkStatusDown    = 0
	LinkStatusUp      = 1
)

type LinkStatus struct {
	Dst2LNK
	Status   int           `json:"status"`
	TTL      time.Duration `json:"-"`
	stored   time.Time
	notified bool
}

func NewLinkStatus(dst, outlink string, status int) *LinkStatus {
	return &LinkStatus{
		Dst2LNK: Dst2LNK{
			DstIP:   dst,
			OutLink: outlink,
		},
		Status: status,
	}
}

func (l *LinkStatus) SetTTL(ttl time.Duration, now ...time.Time) {
	if ttl < 0 {
		ttl = 0
	}
	l.TTL = ttl

	if len(now) == 0 {
		l.stored = time.Now().UTC()
	} else {
		l.stored = now[0].UTC()
	}
}

func (l *LinkStatus) IsExpire(now time.Time) (left time.Duration, ok bool) {
	left = l.TTL - now.UTC().Sub(l.stored)
	return left, left < 0
}
