package models

import (
	"testing"
)

func TestDomainPool(t *testing.T) {
	p, err := GetDomainPool(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(p)
}
