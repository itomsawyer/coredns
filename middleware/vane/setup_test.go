package vane

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
			`vane { 
				upstream_timeout 1s
			}`,
			"perfect",
			false,
		},
		{
			`vane {
				upstream_timeout
			}`,
			"miss upstream_timeout args",
			true,
		},
	}

	for _, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		vane, err := parseVane(c)

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
