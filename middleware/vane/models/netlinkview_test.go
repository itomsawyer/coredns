package models

import (
	"testing"
)

func TestGetNetLink(t *testing.T) {
	query := Values{}
	query.Set("netlink_id", 1)
	nl, err := GetNetLinkView(nil, query, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
