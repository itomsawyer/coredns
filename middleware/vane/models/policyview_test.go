package models

import (
	"testing"
)

func TestPolicyView(t *testing.T) {
	nl, err := GetAllPolicyView(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
