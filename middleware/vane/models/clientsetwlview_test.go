package models

import (
	"testing"
)

func TestGetClientSetWL(t *testing.T) {
	nl, err := GetClientSetWLView(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
