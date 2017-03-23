package engine

import (
	"testing"
)

func TestAddRoute(t *testing.T) {
	rm := RouteMap{}

	ot1 := NewOutLink("out1", "1.1.1.1")
	ot2 := NewOutLink("out2", "1.1.1.2")

	r1 := NewRoute(1, 1, ot1)
	r1.Score = 1
	r1.Priority = 2
	r2 := NewRoute(1, 1, ot2)
	r2.Score = 3
	r2.Priority = 1

	rm.AddRoute(r1)
	rm.AddRoute(r2)

	rs, ok := rm[RouteKey{NetLinkSetID: 1, RouteSetID: 1}]
	if !ok {
		t.Errorf("unexpected not found")
	}

	iter := rs.Iterator()
	t.Log(iter)
	t.Log(iter.Next())
	t.Log(iter.Next())
	t.Log(iter.Next())

	t.Log(rm)
}
