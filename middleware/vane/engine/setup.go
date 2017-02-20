package engine

import (
	"errors"
	"fmt"

	"github.com/miekg/coredns/core/dnsserver"
	"github.com/miekg/coredns/middleware"

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

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		vane.Next = next
		return vane
	})

	return nil
}

func parseVaneEngine(c *caddy.Controller) (vane *VaneEngine, err error) {
	vane = &VaneEngine{
		DBHost: "root:@localhost/iwg",
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
					if len(args) == 0 {
						return nil, c.ArgErr()
					}
					vane.DBHost = args[0]

				default:
					return nil, c.ArgErr()
				}
			}
		}
	}

	if vane.DBHost == "" {
		return nil, ErrNoDBHost
	}

	c.OnFirstStartup(func() error {
		return vane.RegisterDB()
	})

	c.OnStartup(func() error {
		fmt.Println("engine start")
		return vane.Reload()
	})

	return vane, nil
}
