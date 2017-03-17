package models

import (
	"testing"
)

func TestGetSrcView(t *testing.T) {
	nl, err := GetSrcView(nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
