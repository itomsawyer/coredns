package models

import (
	"testing"
)

func TestGetRouteView(t *testing.T) {
	nl, err := GetAllRouteView(nil, nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
