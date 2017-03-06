package nettree

import (
	"net"
	"testing"
)

func TestGetBitFromInt(t *testing.T) {
	if getBit(0x80000000, 0) != 1 {
		t.Errorf("bitK32 failed")
	}

	if getBit(0xf0000000, 0) != 1 {
		t.Errorf("bitK32 failed")
	}

	if getBit(0xf0000000, 3) != 1 {
		t.Errorf("bitK32 failed")
	}

	if getBit(0xf0000000, 4) != 0 {
		t.Errorf("bitK32 failed")
	}

	if getBit(0x00000001, 30) != 0 {
		t.Errorf("bitK32 failed")
	}

	if getBit(0x00000001, 31) != 1 {
		t.Errorf("bitK32 failed")
	}
}

func TestingGetExactly(t *testing.T) {
	nt := new(NetTree)
	nt.Insert(1, 32, 1, 0)
	nt.Insert(7, 32, 7, 0)
	if v := nt.Find(1, 32); v.(int) != 1 {
		t.Errorf("shoud find 1")
	}
	if v := nt.Find(7, 32); v.(int) != 7 {
		t.Errorf("shoud find 7")
	}
}

func TestCovering(t *testing.T) {
	nt := new(NetTree)
	err := nt.Insert(0x0100, 24, 1, 0)
	if err != nil {
		t.Error(err)
	}
	err = nt.Insert(0x0300, 24, 3, 0)
	if err != nil {
		t.Error(err)
	}
	if v := nt.Find(0x0103, 25); v.(int) != 1 {
		t.Errorf("shoud find 1")
	}
	if v := nt.Find(0x0308, 32); v.(int) != 3 {
		t.Errorf("shoud find 3")
	}
}

func TestOverlap(t *testing.T) {
	nt := new(NetTree)
	err := nt.Insert(0x0100, 24, 1, 0)
	if err != nil {
		t.Error(err)
	}
	err = nt.Insert(0x01f0, 28, 2, 0)
	if err != nil {
		t.Error(err)
	}
	if v := nt.Find(0x01ff, 27); v.(int) != 1 {
		t.Errorf("shoud find 1")
	}
	if v := nt.Find(0x01ff, 28); v.(int) != 2 {
		t.Errorf("shoud find 2")
	}
	if v := nt.Find(0x01ff, 32); v.(int) != 2 {
		t.Errorf("shoud find 2")
	}
}

func TestRoot(t *testing.T) {
	nt := new(NetTree)
	err := nt.Insert(0, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}
	err = nt.Insert(0x01f0, 28, 1, 0)
	if err != nil {
		t.Error(err)
	}
	if v := nt.Find(0x01ff, 0); v.(int) != 0 {
		t.Errorf("shoud find 0")
	}
	if v := nt.Find(0x01ff, 24); v.(int) != 0 {
		t.Errorf("shoud find 0")
	}
	if v := nt.Find(0x01ff, 32); v.(int) != 1 {
		t.Errorf("shoud find 1")
	}
}

func TestEmpty(t *testing.T) {
	nt := new(NetTree)

	if v := nt.Find(0x01ff, 1); v != nil {
		t.Errorf("unexpected found")
	}

	if v := nt.Find(0x01ff, 32); v != nil {
		t.Errorf("unexpected found")
	}

	if v := nt.Find(0, 0); v != nil {
		t.Errorf("unexpected found")
	}
}

func TestOverlapByIPNet(t *testing.T) {
	nt := new(NetTree)

	_, cidr, _ := net.ParseCIDR("1.1.1.128/24")
	err := nt.InsertByIPNet(cidr, 0, 0)
	if err != nil {
		t.Error(err)
	}
	_, cidr, _ = net.ParseCIDR("1.1.1.128/28")
	err = nt.InsertByIPNet(cidr, 1, 0)
	if err != nil {
		t.Error(err)
	}

	if v := nt.FindByIP(net.ParseIP("1.1.1.1")); v.(int) != 0 {
		t.Errorf("shoud find 0")
	}
	if v := nt.FindByIP(net.ParseIP("1.1.1.129")); v.(int) != 1 {
		t.Errorf("shoud find 1")
	}
}

func TestOverlapByIPNetWithPrio(t *testing.T) {
	nt := new(NetTree)

	_, cidr, _ := net.ParseCIDR("1.1.0.0/16")
	err := nt.InsertByIPNet(cidr, -1, 1)
	if err != nil {
		t.Error(err)
	}

	_, cidr, _ = net.ParseCIDR("1.1.1.128/24")
	err = nt.InsertByIPNet(cidr, 0, 1)
	if err != nil {
		t.Error(err)
	}

	_, cidr, _ = net.ParseCIDR("1.1.1.128/28")
	err = nt.InsertByIPNet(cidr, 1, 0)
	if err != nil {
		t.Error(err)
	}

	if v := nt.FindByIP(net.ParseIP("1.1.0.1")); v.(int) != -1 {
		t.Errorf("shoud find -1")
	}
	if v := nt.FindByIP(net.ParseIP("1.1.1.1")); v.(int) != 0 {
		t.Errorf("shoud find 0")
	}
	if v := nt.FindByIP(net.ParseIP("1.1.1.129")); v.(int) != 0 {
		t.Errorf("shoud find 0")
	}
}
