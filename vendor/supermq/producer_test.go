package supermq

import (
	"strconv"
	"testing"
)

type IPRoute struct {
	DstIP   string `json:"dst_ip"`
	OutLink string `json:"out_link"`
	Status  bool   `json:"status"`
}

func (i *IPRoute) Output(a int, b string) error {
	return nil
}

func TestMProduer(t *testing.T) {
	mp := NewMProducer()
	if err := mp.AddRx("127.0.0.1:4150"); err != nil {
		t.Fatal(err)
	}
	mp.SetLogger(&IPRoute{}, LogLevelDebug)
	i := 0
	for {
		i++
		i = i % 10
		strI := strconv.Itoa(i)
		mp.MultiPublish("dst2otlnk", [][]byte{
			[]byte(`{"key":"` + strI + `","dst_ip":"127.0.0.` + strI + `","out_link":"HZ-CT"}`)})
	}
}
