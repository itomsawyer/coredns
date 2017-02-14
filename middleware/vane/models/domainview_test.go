package models

import (
	"testing"
)

func TestDomainView(t *testing.T) {
	nl, err := GetAllDomainView(nil, nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
