package vane

import (
	"errors"
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
	v, err := parseVane(c)
	if err != nil {
		return err
	}

	c.OnStartup(v.Init)
	c.OnShutdown(v.Destroy)

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		v.Next = next
		return v
	})

	return nil
}

func parseVane(c *caddy.Controller) (vane *Vane, err error) {
	vane = NewVane()

	for c.Next() {
		if c.Val() == "vane" {
			args := c.RemainingArgs()
			if len(args) > 0 {
				return nil, c.ArgErr()
			}

			for c.NextBlock() {
				switch c.Val() {
				case "keep_cname_chain":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}

					vane.KeepCNAMEChain = parseBool(args[0])

				case "keep_upstream_ecs":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}

					vane.KeepUpstreamECS = parseBool(args[0])

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
						return nil, err
					}

					vane.LogConfigs = append(vane.LogConfigs, lc)

				case "debug":
					args := c.RemainingArgs()
					if len(args) != 1 {
						return nil, c.ArgErr()
					}

					vane.Debug = true

				default:
					return nil, c.ArgErr()
				}
			}
		}
	}

	return vane, nil
}

func parseBool(s string) bool {
	if s == "yes" || s == "on" || s == "true" {
		return true
	}

	return false
}
