package models

import (
	"testing"
)

func TestGetClientSet(t *testing.T) {
	nl, err := GetClientSetView(nil, nil, nil, nil, 0, 10)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
