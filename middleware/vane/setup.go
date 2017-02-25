package vane

import (
	"errors"
	"fmt"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/vane/engine"

	"github.com/mholt/caddy"
)

var (
	ErrNoDBHost = errors.New("invalid db param")
)

func init() {
	caddy.RegisterPlugin("vane", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fmt.Println("setup vane")
	v, err := parseVane(c)
	if err != nil {
		return err
	}

	c.OnStartup(func() error {
		return v.CreateLogger()
	})

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		v.Next = next
		return v
	})

	return nil
}

func parseVane(c *caddy.Controller) (vane *Vane, err error) {
	vane = new(Vane)

	for c.Next() {
		if c.Val() == "vane" {
			args := c.RemainingArgs()
			if len(args) > 0 {
				return nil, c.ArgErr()
			}

			for c.NextBlock() {
				switch c.Val() {
				case "upstream_timeout":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}

					to, err := time.ParseDuration(args[0])
					if err != nil {
						return nil, c.SyntaxErr("time duration")
					}

					vane.UpstreamTimeout = to

				case "log":
					lc, err := engine.ParseLogConfig(c)
					if err != nil {
						fmt.Println("parse log failed")
						return nil, err
					}

					fmt.Println(lc)

					vane.LogConfigs = append(vane.LogConfigs, lc)
				case "debug":
					args := c.RemainingArgs()
					if len(args) != 0 {
						return nil, c.ArgErr()
					}

					vane.Debug = true
				default:
					return nil, c.ArgErr()
				}
			}
		}
	}

	fmt.Println(vane)

	return vane, nil
}
