package engine

import (
	//"fmt"
	//"strings"
	"testing"

	"github.com/mholt/caddy"
)

func TestSetupVane(t *testing.T) {
	tests := []struct {
		input     string
		name      string
		expectErr bool
	}{
		// positive
		{
			`vane_engine { 
				db  root:@localhost/igw
			}`,
			"perfect",
			false,
		},

		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  {
					cache_cap  1000
					unknown_ttl 1s
					enable true
				}
			}`,
			"perfect",
			false,
		},

		{
			`vane_engine {
				db 
			}`,
			"miss db args",
			true,
		},
		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  
			}`,
			"missing lm block",
			true,
		},

		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  {
					cache_cap  1
					unknown_ttl -1s
					enable true
				}
			}`,
			"unknown_ttl too small",
			true,
		},
	}

	for _, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		vane, err := parseVaneEngine(c)

		if test.expectErr {
			if err == nil {
				t.Error(test.name, "expectErr")
				t.Fail()
				return
			}
		} else {
			if err != nil {
				t.Error(test.name, err)
				t.Fail()
				return
			}
		}

		t.Log("vane:", vane)
	}
}
