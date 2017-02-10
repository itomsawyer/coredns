package vane

import (
	"errors"

	"github.com/miekg/coredns/core/dnsserver"
	"github.com/miekg/coredns/middleware"
	"github.com/miekg/coredns/middleware/vane/models"

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
	vane, err := parseVane(c)
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		v := vane
		v.Next = next
		return v
	})

	return nil
}

func parseVane(c *caddy.Controller) (vane *Vane, err error) {
	vane = &Vane{
		DBHost: "root:@localhost/igw",
	}

	for c.Next() {
		if c.Val() == "vane" {
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
		return models.InitDB(vane.DBHost)
	})

	return vane, nil
}
