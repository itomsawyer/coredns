package errors

import (
	"testing"

	"github.com/mholt/caddy"
)

func TestErrorsParse(t *testing.T) {
	tests := []struct {
		inputErrorsRules     string
		shouldErr            bool
		expectedErrorHandler errorHandler
	}{
		{`errors`, false, errorHandler{
			LogFile: "",
		}},
		{`errors errors.txt`, false, errorHandler{
			LogFile: "errors.txt",
		}},
		{`errors visible`, false, errorHandler{
			LogFile: "",
			Debug:   true,
		}},
		{`errors { 
			log visible
		}`, false, errorHandler{
			LogFile: "",
			Debug:   true,
		}},
		{`errors { 
			log visible
			max_backups 3
			max_age 1
			max_size 2
		}`, false, errorHandler{
			LogFile:    "",
			Debug:      true,
			MaxAge:     1,
			MaxSize:    2,
			MaxBackups: 3,
		}},
		{`errors { 
			log errors.txt
			max_backups 3
			max_size 2
		}`, false, errorHandler{
			LogFile:    "errors.txt",
			Debug:      false,
			MaxSize:    2,
			MaxBackups: 3,
		}},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.inputErrorsRules)
		actualErrorsRule, err := errorsParse(c)

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}
		if actualErrorsRule.LogFile != test.expectedErrorHandler.LogFile {
			t.Errorf("Test %d expected LogFile to be %s, but got %s",
				i, test.expectedErrorHandler.LogFile, actualErrorsRule.LogFile)
		}
		if actualErrorsRule.Debug != test.expectedErrorHandler.Debug {
			t.Errorf("Test %d expected Debug to be %v, but got %v",
				i, test.expectedErrorHandler.Debug, actualErrorsRule.Debug)
		}

		if actualErrorsRule.MaxAge != test.expectedErrorHandler.MaxAge {
			t.Errorf("Test %d expected MaxAge to be %v, but got %v",
				i, test.expectedErrorHandler.MaxAge, actualErrorsRule.MaxAge)
		}
		if actualErrorsRule.MaxSize != test.expectedErrorHandler.MaxSize {
			t.Errorf("Test %d expected MaxSize to be %v, but got %v",
				i, test.expectedErrorHandler.MaxSize, actualErrorsRule.MaxSize)
		}
		if actualErrorsRule.MaxBackups != test.expectedErrorHandler.MaxBackups {
			t.Errorf("Test %d expected MaxBackups to be %v, but got %v",
				i, test.expectedErrorHandler.MaxBackups, actualErrorsRule.MaxBackups)
		}
	}
}
