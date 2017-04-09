package engine

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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
	CtlHost    string

	reloadLock sync.Mutex

	ln         net.Listener
	mux        *http.ServeMux
	httpServer *http.Server
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

func (v *VaneEngine) Init() error {
	if err := v.initLogger(); err != nil {
		return err
	}

	go func() {
		var ln net.Listener
		var err error
		var n int

		for {
			ln, err = net.Listen("tcp", v.CtlHost)
			if err != nil {
				if n >= 3 {
					v.Logger.Error("Failed to start vane engine controller listener: %s", err)
				}

				time.Sleep(1 * time.Second)
				n++
				continue
			}
			break
		}

		v.ln = ln
		v.mux = http.NewServeMux()
		v.mux.HandleFunc("/coredns/reload", v.ReloadHandler)

		go func() {
			v.Logger.Info("vane engine controller start to serve")
			v.httpServer = &http.Server{Handler: v.mux}
			v.httpServer.Serve(v.ln)
		}()
	}()

	return nil
}

func (v *VaneEngine) ReloadHandler(w http.ResponseWriter, r *http.Request) {
	err := v.Reload()
	if err != nil {
		fmt.Fprintf(w, "ERR %s\n", err)
	} else {
		fmt.Fprintf(w, "OK\n")
	}
}

func (v *VaneEngine) initLogger() error {
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

	if v.httpServer != nil {
		v.httpServer.Close()
	}
}

func (v *VaneEngine) Reload() error {
	v.reloadLock.Lock()
	defer v.reloadLock.Unlock()

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
