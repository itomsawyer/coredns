package vane

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

					if len(args) >= 2 {
						if e, ok := DBEngines[args[0]]; !ok {
							return nil, fmt.Errorf("DB engine %s does not exists", args[0])
						} else {
							vane.DB = e
						}

						vane.DBHost = args[1]
					} else {
						vane.DB = DBEngines["default"]
						vane.DBHost = args[0]
					}

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
		return vane.DB.Open(vane.DBHost)
	})

	c.OnStartup(func() error {
		return vane.DB.Load()
	})

	return vane, nil
}
