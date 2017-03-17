package supermq

import (
	"strconv"
	//	"strings"
	"testing"
	"time"
)

type IPRoute struct {
	DstIP   string `json:"dst_ip"`
	OutLink string `json:"out_link"`
	Status  bool   `json:"status"`
}

func TestMProduer(t *testing.T) {
	mp := NewMProducer()
	if err := mp.AddRx("127.0.0.1:4150"); err != nil {
		t.Fatal(err)
	}
	i := 0
	for {
		i++
		i = i % 10
		//ta := strings.Repeat("abc", i)
		time.Sleep(time.Second)
		mp.MultiPublish("dst2otlnk", [][]byte{
			[]byte(`{"dst_ip":"127.0.0.` + strconv.Itoa(i) + `","out_link":"changshangguangdian"}`)})
		//mp.MultiPublish("dst2otlnk", [][]byte{[]byte(ta)})
	}
}
