package engine

import (
	"fmt"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/vane/models"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

type VaneEngine struct {
	Next   middleware.Handler
	E      *Engine
	DBName string
	DBHost string
}

func (v VaneEngine) Engine() *Engine {
	return v.E
}

func (v *VaneEngine) Name() string { return "vane_engine" }

func (v *VaneEngine) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	fmt.Println("Enter VaneEngine")
	vctx := context.WithValue(ctx, "vane_engine", v.Engine())
	ret, err := middleware.NextOrFailure(v.Name(), v.Next, vctx, w, r)
	fmt.Println("Leave VaneEngine")
	return ret, err
}

func (v *VaneEngine) RegisterDB() error {
	return models.RegisterDB(v.DBName, "mysql", v.DBHost)
}

func (v *VaneEngine) Reload() error {
	builder := &EngineBuilder{DBName: v.DBName}
	if err := builder.Load(); err != nil {
		return err
	}

	e := new(Engine)
	if err := builder.Build(e); err != nil {
		return err
	}

	v.E = e
	return nil
}
