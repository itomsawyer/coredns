package engine

import (
	"strconv"
	//"sync/atomic"
	"encoding/json"
	"time"

	"github.com/coredns/coredns/middleware/pkg/singleflight"
	"github.com/itomsawyer/llog"

	"supermq"

	"github.com/hashicorp/golang-lru"
	"github.com/mholt/caddy"
	"github.com/nsqio/go-nsq"
)

type LinkManager struct {
	*lru.Cache
	group  *singleflight.Group
	logger *llog.Logger

	LinkStatusTTL  time.Duration
	LinkUnknownTTL time.Duration

	sendHosts []string
	readHosts []string
	SendTopic string
	reader    *supermq.Tconsumer
	sender    *supermq.MProducer
}

func NewLinkManager(cap int) (*LinkManager, error) {
	if cap <= 0 {
		cap = 1000
	}

	c, err := lru.New(cap)
	if err != nil {
		return nil, err
	}

	lm := &LinkManager{
		Cache:          c,
		group:          new(singleflight.Group),
		LinkUnknownTTL: 10 * time.Second,
		LinkStatusTTL:  2 * time.Minute,
	}

	return lm, nil
}

func (m *LinkManager) SetLogger(log *llog.Logger) error {
	m.logger = log

	if m.sender != nil {
		m.sender.SetLogger(log, llogLevel2nsqLogLevel(log.Level()))
	}

	if m.reader != nil {
		m.reader.SetLogger(log, llogLevel2nsqLogLevel(log.Level()))
	}
	return nil
}

func llogLevel2nsqLogLevel(lv int) int {
	if lv <= llog.Ldebug {
		return int(nsq.LogLevelDebug)
	}

	if lv <= llog.Linfo {
		return int(nsq.LogLevelInfo)
	}

	if lv <= llog.Lwarn {
		return int(nsq.LogLevelWarning)
	}

	return int(nsq.LogLevelError)
}

func (m *LinkManager) handlerFunc() supermq.Handle {
	return func(msg *supermq.Message) error {
		ls := &LinkStatus{}
		err := json.Unmarshal(msg.Body, ls)
		if err != nil {
			return nil
		}

		m.logger.Debug("lm handler receive msg:%v", ls)
		ls.SetTTL(m.LinkStatusTTL)
		m.Cache.Add(ls.Dst2LNK, ls)
		return nil
	}
}

func (m *LinkManager) RegisterSender(hosts []string, topic string) error {
	p := supermq.NewMProducer()
	for _, host := range hosts {
		err := p.AddRx(host)
		if err != nil {
			m.logger.Error("Register sender %s to %s error: %s", topic, host, err.Error())
			return err
		}
	}

	m.SendTopic = topic
	m.sendHosts = hosts
	m.sender = p
	return nil
}

func (m *LinkManager) RegisterReader(hosts []string, topic, channel string) error {
	c, err := supermq.NewTconsumer(topic, channel, m.handlerFunc())
	if err != nil {
		m.logger.Error("Register reader %s:%s to %s error: %s", topic, channel, hosts, err.Error())
		return err
	}

	if err := c.ConnectToNSQLookupds(hosts); err != nil {
		m.logger.Error("Register reader %s:%s to %s error: %s", topic, channel, hosts, err.Error())
		return err
	}

	m.readHosts = hosts
	m.reader = c
	return nil
}

func (m *LinkManager) addTaskAndNotify(ls *LinkStatus) error {
	ls.SetTTL(m.LinkUnknownTTL)
	ls.Status = LinkStatusUnknown
	m.Cache.Add(ls.Dst2LNK, ls)

	err := m.registerLink(ls)
	if err != nil {
		return err
	}

	return nil
}

func (m *LinkManager) GetLink(dst, outlink string) (*LinkStatus, bool) {
	ls := NewLinkStatus(dst, outlink, LinkStatusUnknown)

	v, ok := m.Cache.Get(ls.Dst2LNK)
	if !ok {
		m.logger.Debug("link %v not found", ls.Dst2LNK)
		m.addTaskAndNotify(ls)
		return nil, false
	}

	ls = v.(*LinkStatus)
	left, ok := ls.IsExpire(time.Now())
	if ok {
		m.logger.Debug("link %v expired: %s", ls.Dst2LNK, left)
		m.Cache.Remove(ls.Dst2LNK)

		m.addTaskAndNotify(ls)
		return nil, false
	}

	if left*2 < ls.TTL && ls.notified != true {
		m.logger.Debug("link %v almost expired: %s", ls.Dst2LNK, left)
		err := m.registerLink(ls)
		if err == nil {
			ls.notified = true
		}
	}

	m.logger.Debug("link found %v", ls)
	return ls, true
}

func (m *LinkManager) registerLink(ls *LinkStatus) error {
	if m.sender == nil {
		return nil
	}

	_, err := m.group.Do(ls.Dst2LNK.DstIP+":"+ls.Dst2LNK.OutLink, func() (interface{}, error) {
		data, err := json.Marshal(ls)
		if err != nil {
			m.logger.Warn("register link status pack error: %s", err)
			return nil, err
		}

		if err := m.sender.MultiPublish(m.SendTopic, [][]byte{data}); err != nil {
			m.logger.Warn("register link status send error: %s", err)
			return nil, err
		}

		m.logger.Debug("register link: %s success", ls.Dst2LNK)
		return nil, nil
	})

	return err
}

func (m *LinkManager) Stop() {
	if m.sender != nil {
		m.sender.Stop()
		m.logger.Info("stop link manager sender")
	}

	if m.reader != nil {
		m.reader.Stop()
		m.logger.Info("stop link manager reader")
	}

	if m.logger != nil {
		m.logger.Close()
	}
}

type LinkManagerConfig struct {
	Enable         bool
	Cap            int
	SendTopic      string
	ReadTopic      string
	ReadChannel    string
	sendHosts      []string
	readHosts      []string
	LinkStatusTTL  time.Duration
	LinkUnknownTTL time.Duration
	LogConfig      *llog.Config
}

func (c *LinkManagerConfig) CreateLinkManager() (*LinkManager, error) {
	lm, err := NewLinkManager(c.Cap)
	if err != nil {
		return nil, err
	}

	lm.LinkUnknownTTL = c.LinkUnknownTTL
	lm.LinkStatusTTL = c.LinkStatusTTL
	if logger, err := CreateLogger(c.LogConfig); err != nil {
		return nil, err
	} else {
		if err := lm.SetLogger(logger); err != nil {
			return nil, err
		}
	}

	lm.logger.Info("Register sender %s to %v", c.SendTopic, c.sendHosts)
	if err := lm.RegisterSender(c.sendHosts, c.SendTopic); err != nil {
		return nil, err
	}

	lm.logger.Info("Register read %s:%s to %v", c.ReadTopic, c.ReadChannel, c.readHosts)
	if err := lm.RegisterReader(c.readHosts, c.ReadTopic, c.ReadChannel); err != nil {
		return nil, err
	}

	return lm, nil
}

func NewLinkManagerConfig() *LinkManagerConfig {
	return &LinkManagerConfig{
		Enable:         true,
		Cap:            100,
		LinkUnknownTTL: 10 * time.Second,
		LinkStatusTTL:  2 * time.Minute,
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

		case "link_ttl":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			if lmconfig.LinkStatusTTL, err = time.ParseDuration(args[0]); err != nil {
				return nil, c.SyntaxErr(err.Error())
			}

			if lmconfig.LinkStatusTTL <= 0 {
				return nil, c.SyntaxErr("unknown_ttl should be greater than 0")
			}

		case "send_topic":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			lmconfig.SendTopic = args[0]

		case "read_topic":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			lmconfig.ReadTopic = args[0]

		case "read_channel":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			lmconfig.ReadChannel = args[0]

		case "send_hosts":
			args := c.RemainingArgs()
			if len(args) == 0 {
				return nil, c.ArgErr()
			}

			lmconfig.sendHosts = args

		case "read_hosts":
			args := c.RemainingArgs()
			if len(args) == 0 {
				return nil, c.ArgErr()
			}

			lmconfig.readHosts = args

		case "log":
			lc, err := ParseLogConfig(c)
			if err != nil {
				return nil, err
			}

			lmconfig.LogConfig = lc
		default:
			return nil, c.Err("directive " + c.Val() + " is unknown")
		}
	}

	if lmconfig.SendTopic == "" || lmconfig.ReadTopic == "" || lmconfig.ReadChannel == "" {
		return nil, c.Err("send_topic read_topic and read_channel must be set")
	}

	if len(lmconfig.sendHosts) == 0 || len(lmconfig.readHosts) == 0 {
		return nil, c.Err("send_hosts and read_hosts  must be set")
	}

	return lmconfig, nil
}
