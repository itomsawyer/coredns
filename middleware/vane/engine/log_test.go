package engine

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestSetupLog(t *testing.T) {
	tests := []struct {
		input     string
		name      string
		expectErr bool
	}{
		// positive
		{
			`log { 
				type console
				level info
			}`,
			"perfect",
			false,
		},

		{
			`log { 
				type nonexists
				level info
			}`,
			"wrong log type",
			true,
		},
	}

	for _, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		c.Next()
		log, err := ParseLogConfig(c)

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

		t.Log("log config: ", log)
	}
}
