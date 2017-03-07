package engine

import (
	"net"
	"strconv"
	//"sync/atomic"
	"time"

	"github.com/coredns/coredns/middleware/pkg/singleflight"

	"github.com/hashicorp/golang-lru"
	"github.com/mholt/caddy"
)

const (
	LinkStatusUnknown = iota
	LinkStatusUp
	LinkStatusDown
)

type LinkManager struct {
	*lru.Cache
	group *singleflight.Group

	LinkUnknownTTL time.Duration
}

func NewLinkManager(cap int) (*LinkManager, error) {
	if cap <= 0 {
		cap = 1000
	}

	c, err := lru.New(cap)
	if err != nil {
		return nil, err
	}

	return &LinkManager{
		Cache:          c,
		group:          new(singleflight.Group),
		LinkUnknownTTL: 3 * time.Second,
	}, nil
}

func (m *LinkManager) Launch() {

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

func (l *LinkStatus) MarkNotify() {
	l.notified = true
}

func (l *LinkStatus) IsNofitied() bool {
	return l.notified
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
		ls.SetTTL(m.LinkUnknownTTL)
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

	if left*2 < ls.TTL && ls.IsNofitied() {
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

type LinkManagerConfig struct {
	Enable         bool
	Cap            int
	LinkUnknownTTL time.Duration
}

func NewLinkManagerConfig() *LinkManagerConfig {
	return &LinkManagerConfig{
		Enable:         true,
		Cap:            100,
		LinkUnknownTTL: 3 * time.Second,
	}
}

func ParseLinkManagerConfig(c *caddy.Controller) (*LinkManagerConfig, error) {
	var err error

	if c.Val() != "lm" {
		return nil, c.SyntaxErr("lm")
	}
	args := c.RemainingArgs()

	//jump over log
	c.Next()
	for range args {
		//jump over RemainingArgs
		c.Next()
	}

	if c.Val() != "{" {
		return nil, c.SyntaxErr("expect {")
	}

	//Config block nest anoter block
	c.IncrNest()

	lmconfig := NewLinkManagerConfig()
	for c.NextBlock() {
		switch c.Val() {
		case "enable":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}

			if args[0] == "yes" || args[0] == "on" {
				lmconfig.Enable = true
			}

		case "cache_cap":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			if lmconfig.Cap, err = strconv.Atoi(args[0]); err != nil {
				return nil, c.SyntaxErr(err.Error())
			}

			if lmconfig.Cap <= 0 {
				return nil, c.SyntaxErr("cache_cap should be greater than 0")
			}

		case "unknown_ttl":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			if lmconfig.LinkUnknownTTL, err = time.ParseDuration(args[0]); err != nil {
				return nil, c.SyntaxErr(err.Error())
			}

			if lmconfig.LinkUnknownTTL <= 0 {
				return nil, c.SyntaxErr("unknown_ttl should be greater than 0")
			}

		default:
			return nil, c.ArgErr()
		}
	}

	return lmconfig, nil
}
