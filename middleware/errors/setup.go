package errors

import (
	"io"
	"log"
	"os"
	"strconv"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/hashicorp/go-syslog"
	"github.com/mholt/caddy"
	"github.com/natefinch/lumberjack"
)

func init() {
	caddy.RegisterPlugin("errors", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	handler, err := errorsParse(c)
	if err != nil {
		return middleware.Error("errors", err)
	}

	var writer io.Writer
	var logFile bool

	switch handler.LogFile {
	case "visible":
		handler.Debug = true
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "syslog":
		writer, err = gsyslog.NewLogger(gsyslog.LOG_ERR, "LOCAL0", "coredns")
		if err != nil {
			return middleware.Error("errors", err)
		}
	default:
		if handler.LogFile == "" {
			writer = os.Stderr // default
			break
		}

		var file *os.File
		file, err = os.OpenFile(handler.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return middleware.Error("errors", err)
		}
		file.Close()
		logFile = true
	}

	handler.Log = log.New(writer, "", 0)
	if logFile {
		lj := &lumberjack.Logger{
			Filename:   handler.LogFile,
			MaxSize:    handler.MaxSize,    // megabytes after which new file is created
			MaxBackups: handler.MaxBackups, // number of backups
			MaxAge:     handler.MaxAge,     //days
			LocalTime:  true,
		}
		handler.Log.SetOutput(lj)
		handler.LogCloser = lj
	}

	c.OnShutdown(handler.Close)

	dnsserver.GetConfig(c).AddMiddleware(func(next middleware.Handler) middleware.Handler {
		handler.Next = next
		return handler
	})

	return nil
}

func errorsParse(c *caddy.Controller) (errorHandler, error) {
	handler := errorHandler{}

	for c.Next() {
		args := c.RemainingArgs()

		if len(args) > 0 {
			if args[0] == "visible" {
				handler.Debug = true
			} else {
				handler.LogFile = c.Val()
			}
		}

		for c.NextBlock() {
			switch c.Val() {
			case "max_age":
				maxAge := c.RemainingArgs()
				if len(maxAge) == 0 {
					return handler, c.ArgErr()
				}
				v, err := strconv.Atoi(maxAge[0])
				if err != nil {
					return handler, c.SyntaxErr("max_age parse error" + err.Error())
				}
				handler.MaxAge = v
			case "max_backups":
				maxBackups := c.RemainingArgs()
				if len(maxBackups) == 0 {
					return handler, c.ArgErr()
				}
				v, err := strconv.Atoi(maxBackups[0])
				if err != nil {
					return handler, c.SyntaxErr("max_backups parse error" + err.Error())
				}
				handler.MaxBackups = v
			case "max_size":
				maxSize := c.RemainingArgs()
				if len(maxSize) == 0 {
					return handler, c.ArgErr()
				}
				v, err := strconv.Atoi(maxSize[0])
				if err != nil {
					return handler, c.SyntaxErr("max_size parse error" + err.Error())
				}
				handler.MaxSize = v
			case "log":
				where := c.RemainingArgs()
				if len(where) == 0 {
					return handler, c.ArgErr()
				}

				if where[0] == "visible" {
					handler.Debug = true
				} else {
					handler.LogFile = where[0]
				}

			default:
				return handler, c.ArgErr()
			}
		}
	}

	return handler, nil
}
