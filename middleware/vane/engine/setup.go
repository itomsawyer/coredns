package engine

import (
	"errors"
	"net"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/mholt/caddy"
)

var (
	ErrNoDBHost = errors.New("invalid db param")
)

func init() {
	caddy.RegisterPlugin("vane_engine", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	vane, err := parseVaneEngine(c)
	if err != nil {
		return err
	}

	c.OnFirstStartup(func() error {
		return vane.RegisterDB()
	})

	c.OnStartup(func() (err error) {
		err = vane.Init()
		return err
	})

	c.OnStartup(func() (err error) {
		err = vane.Reload()
		if err == nil {
			vane.Logger.Info("vane engine start success")
		} else {
			vane.Logger.Info("vane engine start fail")
		}
		return err
	})

	c.OnShutdown(func() error {
		vane.Stop()
		return nil
	})

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		vane.Next = next
		return vane
	})

	return nil
}

func parseVaneEngine(c *caddy.Controller) (vane *VaneEngine, err error) {
	vane = &VaneEngine{
		DBHost:  "root:@localhost/iwg",
		CtlHost: "127.0.0.1:9053",
		cancel:  make(chan struct{}, 1),
	}

	for c.Next() {
		if c.Val() == "vane_engine" {
			args := c.RemainingArgs()
			if len(args) > 0 {
				return nil, c.ArgErr()
			}

			for c.NextBlock() {
				switch c.Val() {
				case "db":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}

					vane.DBHost = args[0]

				case "ctl_host":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}
					vane.CtlHost = args[0]
					_, _, err := net.SplitHostPort(vane.CtlHost)
					if err != nil {
						return nil, c.SyntaxErr(err.Error())
					}

				case "log":
					lc, err := ParseLogConfig(c)
					if err != nil {
						return nil, err
					}

					vane.LogConfigs = append(vane.LogConfigs, lc)
				case "lm":
					lmc, err := ParseLinkManagerConfig(c)
					if err != nil {
						return nil, err
					}

					vane.LMConfig = lmc
				default:
					return nil, c.ArgErr()
				}
			}
		}
	}

	if vane.DBHost == "" {
		return nil, ErrNoDBHost
	}

	return vane, nil
}
