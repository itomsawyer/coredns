package dmtree

import (
	"strings"

	"github.com/miekg/dns"
)

type DmTree struct {
	Next      map[string]*DmTree
	Value     interface{}
	WildValue interface{}
}

func (t *DmTree) Find(domain string) (interface{}, bool) {
	var (
		i     int
		v     interface{}
		found *DmTree
		next  *DmTree
	)

	if dns.IsFqdn(domain) {
		domain = domain[:len(domain)-1]
	}

	if len(domain) == 0 {
		return t.Value, true
	}

	tokens := strings.Split(domain, ".")
	next = t
	for i = len(tokens) - 1; i >= 0; i-- {
		if next.WildValue != nil {
			found = next
		}

		key := tokens[i]
		if next.Next == nil {
			break
		}

		node, ok := next.Next[key]
		if !ok {
			break
		}

		next = node
	}

	if i < 0 {
		found = next
		v = found.Value
		return v, true
	}

	if found != nil {
		v = found.WildValue
		return v, true
	}

	return nil, false
}

func (t *DmTree) Insert(domain string, v interface{}) {
	var (
		node *DmTree
	)

	if dns.IsFqdn(domain) {
		domain = domain[:len(domain)-1]
	}

	if len(domain) == 0 {
		t.WildValue = v
		return
	}

	if len(domain) >= 2 && domain[:2] == "*." {
		domain = domain[1:]
	}

	tokens := strings.Split(domain, ".")

	node = t
	for i := len(tokens) - 1; i >= 0; i-- {
		key := tokens[i]
		if key == "" {
			node.WildValue = v
			break
		}

		if node.Next == nil {
			node.Next = make(map[string]*DmTree, 4)
		}

		if n, ok := node.Next[key]; ok {
			node = n
		} else {
			n = new(DmTree)
			node.Next[key] = n
			node = n
		}
	}

	node.Value = v
	return
}
