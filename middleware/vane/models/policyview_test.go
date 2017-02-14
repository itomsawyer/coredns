package models

import (
	"testing"
)

func TestPolicyView(t *testing.T) {
	nl, err := GetAllPolicyView(nil, nil, nil, nil, nil, 0, 1)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	t.Log(nl)
}

func TestPolicySlice(t *testing.T) {
	var (
		hosts [][]string
		p     PolicySlice
	)

	p = PolicySlice{
		&PolicyView{Priority: 1, Addr: "1.1.1.1"},
	}

	hosts = p.Hosts()
	t.Log("singel hosts", hosts)

	p = PolicySlice{
		&PolicyView{Priority: 1, Addr: "1.1.1.1"},
		&PolicyView{Priority: 2, Addr: "1.1.2.1"},
		&PolicyView{Priority: 3, Addr: "1.1.3.1"},
	}

	hosts = p.Hosts()
	t.Log("multiple hosts", hosts)

	p = PolicySlice{
		&PolicyView{Priority: 1, Addr: "1.1.1.1"},
		&PolicyView{Priority: 1, Addr: "1.1.1.2"},
		&PolicyView{Priority: 2, Addr: "1.1.2.1"},
		&PolicyView{Priority: 2, Addr: "1.1.2.2"},
		&PolicyView{Priority: 3, Addr: "1.1.3.1"},
	}

	hosts = p.Hosts()
	t.Log("multiple hosts", hosts)

}
