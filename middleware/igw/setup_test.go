package igw

import (
	//"fmt"
	//"strings"
	"testing"

	"github.com/mholt/caddy"
)

func TestSetupDemo(t *testing.T) {
	tests := []struct {
		input string
	}{
		// positive
		{
			`igw`,
		},
	}

	for _, test := range tests {
		caddy.NewTestController("dns", test.input)
	}
}
