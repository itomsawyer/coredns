package engine

import (
	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/vane/models"

	"github.com/astaxie/beego/logs"
	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type VaneEngine struct {
	Next       middleware.Handler
	E          *Engine
	DBName     string
	DBHost     string
	Logger     *logs.BeeLogger
	LogConfigs []*LogConfig
	LMConfig   *LinkManagerConfig
}

func (v VaneEngine) Engine() *Engine {
	return v.E
}

func (v *VaneEngine) Name() string { return "vane_engine" }

func (v *VaneEngine) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	vctx := context.WithValue(ctx, "vane_engine", v.Engine())
	ret, err := middleware.NextOrFailure(v.Name(), v.Next, vctx, w, r)
	return ret, err
}

func (v *VaneEngine) RegisterDB() error {
	return models.RegisterDB(v.DBName, "mysql", v.DBHost)
}

func (v *VaneEngine) InitLogger() error {
	l, err := CreateLogger(v.LogConfigs)
	if err != nil {
		return err
	}

	v.Logger = l
	return nil
}

func (v *VaneEngine) Stop() {
	v.Logger.Info("vane engine stop")
	if v.E != nil && v.E.LinkManager != nil {
		v.E.LinkManager.Stop()
	}
}

func (v *VaneEngine) Reload() error {
	v.Logger.Info("vane engine load")
	builder := &EngineBuilder{DBName: v.DBName, Logger: v.Logger}
	if err := builder.Load(); err != nil {
		return err
	}

	e := new(Engine)
	if err := builder.Build(e); err != nil {
		return err
	}

	if v.LMConfig != nil && v.LMConfig.Enable {
		if lm, err := v.LMConfig.CreateLinkManager(); err != nil {
			return err
		} else {
			e.LinkManager = lm
		}
	}

	v.E = e
	return nil
}
