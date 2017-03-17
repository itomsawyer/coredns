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
			`
			vane_engine {
				log {
					type console	
					level info
				}	
			}
			vane { 
				upstream_timeout 1s
				log {
					type console	
					level error
				}
			}`,
			"perfect",
			false,
		},

		{
			`
			vane {
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
			continue
		} else {
			if err != nil {
				t.Error(test.name, err)
				t.Fail()
				return
			}
		}

		t.Log("vane:", vane)
		if len(vane.LogConfigs) != 0 {
			for _, c := range vane.LogConfigs {
				t.Log("vane log config:", c)
			}
		}
	}
}
