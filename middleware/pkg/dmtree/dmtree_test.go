package dmtree

import (
	"testing"
)

func TestDmTreeRoot(t *testing.T) {
	var (
		ok      bool
		v       interface{}
		domains = []string{"a.", ".", "a", ""}
	)

	dt := new(DmTree)
	err := dt.Insert(".", 1)
	if err != nil {
		t.Errorf("unexpected error")
		return
	}

	err = dt.Insert("", 2)
	if err == nil {
		t.Errorf("unexpected no error")
		return
	}

	for _, d := range domains {
		v, ok = dt.Find("a.")
		if !ok {
			t.Errorf("unexpected not found for %s", d)
		}

		if v == nil || v.(int) != 1 {
			t.Errorf("unexpected value for %s", d)
		}
	}
}

func TestDmTreeForceInsert(t *testing.T) {
	var (
		ok bool
		v  interface{}
	)

	dt := new(DmTree)
	dt.Insert(".", 1)
	dt.ForceInsert("", 2)

	v, ok = dt.Find("a.")
	if !ok {
		t.Errorf("unexpected not found")
		return
	}

	if v == nil || v.(int) != 2 {
		t.Errorf("unexpected value")
		return
	}
}

func TestDmTreeExactly(t *testing.T) {
	var (
		v  interface{}
		ok bool
	)
	dt := new(DmTree)
	dt.Insert("qq.com", 1)
	dt.Insert("baidu.com", 2)

	v, ok = dt.Find("qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("baidu.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 2 {
		t.Errorf("expect value 2 get", v)
	}

	v, ok = dt.Find("long.static.qq.com")
	if ok {
		t.Errorf("unexpected found")
	}

	return
}

func TestDmTreeCovering(t *testing.T) {
	var (
		v  interface{}
		ok bool
	)
	dt := new(DmTree)
	dt.Insert("qq.com", 1)
	dt.Insert("static.qq.com", 2)

	v, ok = dt.Find("qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("static.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 2 {
		t.Errorf("expect value 2 get", v)
	}

	v, ok = dt.Find("long.static.qq.com")
	if ok {
		t.Errorf("unexpected found")
	}

	return
}

func TestDmTreeWild(t *testing.T) {
	var (
		v  interface{}
		ok bool
	)

	dt := new(DmTree)
	dt.Insert(".", 0)
	dt.Insert(".qq.com", 1)
	dt.Insert("static.qq.com", 2)

	v, ok = dt.Find("www.baidu.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 0 {
		t.Errorf("expect value 0 get", v)
	}

	v, ok = dt.Find("qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 0 {
		t.Errorf("expect value 0 get %d", v)
	}

	v, ok = dt.Find("a.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("a.a.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("static.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 2 {
		t.Errorf("expect value 2 get", v)
	}

	//XXX Warning: THIS feature is not the same with dns wildcard domain protocol.
	//             BUT this is what we want.
	v, ok = dt.Find("a.static.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	return
}

func TestDmTreeWild2(t *testing.T) {
	var (
		v  interface{}
		ok bool
	)

	dt := new(DmTree)
	dt.Insert("*.qq.com", 1)
	dt.Insert("static.qq.com", 2)

	v, ok = dt.Find("www.baidu.com")
	if ok {
		t.Errorf("unexpected found")
	}

	v, ok = dt.Find("qq.com")
	if ok {
		t.Errorf("unexpected found")
	}

	v, ok = dt.Find("a.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("a.a.qq.com")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 1 {
		t.Errorf("expect value 1 get", v)
	}

	v, ok = dt.Find("static.qq.com.")
	if !ok {
		t.Errorf("unexpected not found")
	}
	if v == nil || v.(int) != 2 {
		t.Errorf("expect value 2 get", v)
	}

	return
}

func TestDmTreeNilValue(t *testing.T) {
	var ok bool

	dt := new(DmTree)
	dt.Insert(".qq.com", nil)

	_, ok = dt.Find("a.qq.com")
	if ok {
		t.Errorf("unexpected found")
	}

	_, ok = dt.Find("com")
	if ok {
		t.Errorf("unexpected found")
	}

	_, ok = dt.Find("")
	if ok {
		t.Errorf("unexpected found")
	}

	_, ok = dt.Find("static.qq.com.")
	if ok {
		t.Errorf("unexpected found")
	}
}
