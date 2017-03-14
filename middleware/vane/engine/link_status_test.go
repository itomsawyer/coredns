package engine

import (
	"encoding/json"
	"testing"
)

func TestLinkStatusUnmarshal(t *testing.T) {
	var (
		d   Dst2LNK
		ls  LinkStatus
		err error
	)

	ed := Dst2LNK{DstIP: "1.1.1.1", OutLink: "ot1"}
	err = json.Unmarshal([]byte(`{"Dst_ip":"1.1.1.1", "Out_link":"ot1", "Status":1}`), &d)
	if err != nil {
		t.Error(err)
	} else if ed != d {
		t.Errorf("unexpected value expect %v, get %v", ed, d)
	}

	err = json.Unmarshal([]byte(`{"Dst_ip":"1.1.1.1", "Out_link":"ot1", "Status":1}`), &ls)
	if err != nil {
		t.Error(err)
	} else if ls.Dst2LNK != ed || ls.Status != 1 {
		t.Errorf("unexpected value expect %v, get %v", ed, ls)
	}
}
