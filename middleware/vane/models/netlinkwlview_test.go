package models

import (
	"testing"
)

func TestGetNetLinkWL(t *testing.T) {
	query := Values{}
	query.Set("netlink_id", 1)
	nl, err := GetAllNetLinkWLView(nil, query, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}
