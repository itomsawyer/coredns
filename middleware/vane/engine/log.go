package engine

import (
	"encoding/json"
	"errors"

	"github.com/astaxie/beego/logs"
	"github.com/mholt/caddy"
)

var (
	supportLogTypeList = logs.AdapterFile

	supportLevels = map[string]int{
		"error": logs.LevelError,
		"warn":  logs.LevelWarn,
		"info":  logs.LevelInfo,
		"debug": logs.LevelDebug,
	}
)

type LogConfig struct {
	adapter  string
	FileName string
	Level    int
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		adapter: logs.AdapterFile,
		Level:   logs.LevelInfo,
	}
}

func (lc *LogConfig) ApplyTo(l *logs.BeeLogger) error {
	config, err := lc.Marshal()
	if err != nil {
		return err
	}

	return l.SetLogger(lc.adapter, config)
}

func (lc *LogConfig) Marshal() (string, error) {
	b, err := json.Marshal(lc)
	return string(b), err
}

func supportLogType(t string) bool {
	switch t {
	case logs.AdapterFile:
		return true
	default:
		return false
	}
}

func logLevel(l string) (int, error) {
	n, ok := supportLevels[l]
	if !ok {
		return 0, errors.New("unsupported log level:" + l)
	}

	return n, nil
}

func ParseLogConfig(c *caddy.Controller) (*LogConfig, error) {
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
		return nil, c.ArgErr()
	}

	//Config block nest anoter block (log block)
	c.IncrNest()

	lc := NewLogConfig()
	for c.NextBlock() {
		switch c.Val() {
		case "type":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}

			if !supportLogType(args[0]) {
				return nil, c.SyntaxErr(supportLogTypeList)
			}

			lc.adapter = args[0]

		case "filename":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}

			lc.FileName = args[0]
		case "level":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}

			n, err := logLevel(args[0])
			if err != nil {
				return nil, err
			}

			lc.Level = n
		default:
			return nil, c.ArgErr()
		}
	}

	return lc, nil
}
