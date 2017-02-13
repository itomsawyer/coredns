package models

import (
	"testing"
)

func TestGetClientSet(t *testing.T) {
	nl, err := GetAllClientSetView(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
