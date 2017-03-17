package models

import (
	"testing"
)

func TestGetDstView(t *testing.T) {
	nl, err := GetDstView(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
