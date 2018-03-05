package engine

import (
	"strconv"

	"github.com/itomsawyer/llog"
	"github.com/mholt/caddy"
)

func NewLogConfig() *llog.Config {
	return &llog.Config{
		OutputFile: "stdout",
		Level:      "info",
	}
}

func CreateLogger(lc *llog.Config) (*llog.Logger, error) {
	return llog.New(*lc, 0)
}

func ParseLogConfig(c *caddy.Controller) (*llog.Config, error) {
	if c.Val() != "log" {
		return nil, c.SyntaxErr("log")
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

	//Config block nest anoter block (log block)
	c.IncrNest()

	lc := NewLogConfig()
	for c.NextBlock() {
		switch c.Val() {
		case "filename":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			lc.OutputFile = args[0]

		case "level":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			lc.Level = args[0]

		case "max_age":
			maxAge := c.RemainingArgs()
			if len(maxAge) == 0 {
				return nil, c.ArgErr()
			}
			v, err := strconv.Atoi(maxAge[0])
			if err != nil {
				return nil, c.SyntaxErr("max_age parse error" + err.Error())
			}
			lc.MaxAge = v

		case "max_backups":
			maxBackups := c.RemainingArgs()
			if len(maxBackups) == 0 {
				return nil, c.ArgErr()
			}
			v, err := strconv.Atoi(maxBackups[0])
			if err != nil {
				return nil, c.SyntaxErr("max_backups parse error" + err.Error())
			}
			lc.MaxBackups = v

		case "max_size":
			maxSize := c.RemainingArgs()
			if len(maxSize) == 0 {
				return nil, c.ArgErr()
			}
			v, err := strconv.Atoi(maxSize[0])
			if err != nil {
				return nil, c.SyntaxErr("max_size parse error" + err.Error())
			}
			lc.MaxSize = v

		default:
			return nil, c.ArgErr()
		}
	}

	return lc, nil
}
