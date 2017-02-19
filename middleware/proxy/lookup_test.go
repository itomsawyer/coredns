package proxy

import (
	"testing"
)

func TestNewLookup(t *testing.T) {
	hosts := []string{"8.8.8.8", "114.114.114.114"}
	proxy := NewLookup(hosts)
	t.Log(proxy)

	h := [][]string{hosts}
	proxy = NewLookup2(h)
	t.Log(proxy)
}
