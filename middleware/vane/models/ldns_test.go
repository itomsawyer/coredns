package models

import (
	"testing"
)

func TestLDNS(t *testing.T) {
	nl, err := GetLDNS(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
