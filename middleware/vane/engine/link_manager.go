package engine

import (
	"net"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware/pkg/singleflight"

	"github.com/hashicorp/golang-lru"
)

const (
	LinkStatusUnknown = iota
	LinkStatusUp
	LinkStatusDown
)

type LinkManager struct {
	*lru.Cache
	group *singleflight.Group

	LinkUnknownTimeout time.Duration
}

func NewLinkManager(cap int) (*LinkManager, error) {
	c, err := lru.New(cap)
	if err != nil {
		return nil, err
	}

	return &LinkManager{
		Cache: c,
		group: new(singleflight.Group),
	}, nil
}

type LinkStatus struct {
	Dst       net.IP
	OutLinkID int
	Status    int
	TTL       time.Duration
	stored    time.Time
	notified  bool
}

func (l *LinkStatus) Key() string {
	if l.Dst == nil {
		return "nil/" + strconv.Itoa(l.OutLinkID)
	}

	return l.Dst.String() + "/" + strconv.Itoa(l.OutLinkID)
}

func (l *LinkStatus) SetTTL(ttl time.Duration) {
	if ttl < 0 {
		ttl = 0
	}
	l.TTL = ttl

	ls.stored = time.Now().UTC()
}

func (l *LinkStatus) MarkNotify() {
	return l.notified = true
}

func (l *LinkStatus) NeedNotify() bool {
	return l.notified == false
}

func (l *LinkStatus) IsExpire(now time.Time) (left time.Duration, ok bool) {
	left = l.TTL - now.UTC().Sub(l.stored)
	return left, left < 0
}

func (m *LinkManager) GetLink(dst net.IP, outLinkID int) (*LinkStatus, bool) {
	if dst == nil {
		return nil, false
	}

	ls := &LinkStatus{Dst: dst, OutLinkID: outLinkID, Status: LinkStatusUnknown}
	key := ls.Key()

	v, ok := m.Cache.Get(key)
	if !ok {
		ls.SetTTL(m.LinkUnknownTimeout)
		m.Cache.Add(key, ls)
		m.registerLink(key, ls)
		return nil, false
	}

	ls = v.(*LinkStatus)
	left, ok := ls.IsExpire(time.Now())
	if ok {
		m.Cache.Remove(key)
		return nil, false
	}

	if left*2 < ls.TTL && ls.NeedNotify(){
		err := m.registerLink(key, ls)
		if err == nil {
			ls.MarkNotify()
		}
	}

	return ls, true
}

func (m *LinkManager) registerLink(key string, ls *LinkStatus) error {
	_, err := m.group.Do(key, func() (interface{}, error) {
		//TODO register link status to backend

		return nil, nil
	})

	return err
}
