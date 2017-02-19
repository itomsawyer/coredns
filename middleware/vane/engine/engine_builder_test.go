package engine

import (
	"testing"
)

func TestBuilderLoad(t *testing.T) {
	e := new(EngineBuilder)
	e.DBName = "default"

	e.Load()
	t.Log(e.PolicyView)
}

func TestBuildEngine(t *testing.T) {
	b := new(EngineBuilder)
	b.DBName = "default"

	b.Load()
	t.Log(b.PolicyView)

	e := new(Engine)
	err := b.Build(e)
	if err != nil {
		t.Error(err)
	}

	t.Log(e)
}
