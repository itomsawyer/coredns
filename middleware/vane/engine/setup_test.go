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
			"ok",
			false,
		},

		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  {
					cache_cap  1000
					unknown_ttl 1s
					enable true
					read_hosts 127.0.0.1
					send_hosts 127.0.0.1
					send_topic send
					read_topic read
					read_channel read
				}
			}`,
			"ok with lm",
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
					read_hosts 127.0.0.1
					send_hosts 127.0.0.1
					send_topic send
					read_topic read
					read_channel read
				}
			}`,
			"unknown_ttl too small",
			true,
		},
		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  {
					cache_cap  1
					unknown_ttl 2s
					enable true
					read_hosts 127.0.0.1
					send_topic send
					read_topic read
					read_channel read
				}
			}`,
			"send_hosts must be set ",
			true,
		},

		{
			`vane_engine { 
				db  root:@localhost/igw
				lm  {
					cache_cap  1
					unknown_ttl 2s
					enable true
					read_hosts 127.0.0.1
					send_hosts 127.0.0.1
					read_topic read
					read_channel read
				}
			}`,
			"send_topic must be set",
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
		t.Log(vane)
	}
}
