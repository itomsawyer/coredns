package log

import (
	"testing"

	"github.com/coredns/coredns/middleware/pkg/response"

	"github.com/mholt/caddy"
)

func TestLogParse(t *testing.T) {
	tests := []struct {
		inputLogRules    string
		shouldErr        bool
		expectedLogRules []Rule
	}{
		{`log`, false, []Rule{{
			NameScope:  ".",
			OutputFile: DefaultLogFilename,
			Format:     DefaultLogFormat,
		}}},
		{`log log.txt`, false, []Rule{{
			NameScope:  ".",
			OutputFile: "log.txt",
			Format:     DefaultLogFormat,
		}}},
		{`log example.org log.txt`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     DefaultLogFormat,
		}}},
		{`log example.org. stdout`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "stdout",
			Format:     DefaultLogFormat,
		}}},
		{`log example.org log.txt {common}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     CommonLogFormat,
		}}},
		{`log example.org accesslog.txt {combined}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "accesslog.txt",
			Format:     CombinedLogFormat,
		}}},
		{`log example.org. log.txt
			  log example.net accesslog.txt {combined}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     DefaultLogFormat,
		}, {
			NameScope:  "example.net.",
			OutputFile: "accesslog.txt",
			Format:     CombinedLogFormat,
		}}},
		{`log example.org stdout {host}
			  log example.org log.txt {when}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "stdout",
			Format:     "{host}",
		}, {
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     "{when}",
		}}},

		{`log example.org log.txt {
				class all
			}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     CommonLogFormat,
			Class:      response.All,
		}}},
		{`log example.org log.txt {
				class all
				max_age 1
				max_size 2
				max_backups 3
			}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     CommonLogFormat,
			Class:      response.All,
			MaxAge:     1,
			MaxSize:    2,
			MaxBackups: 3,
		}}},
		{`log example.org log.txt {
			class denial
		}`, false, []Rule{{
			NameScope:  "example.org.",
			OutputFile: "log.txt",
			Format:     CommonLogFormat,
			Class:      response.Denial,
		}}},
		{`log {
			class denial
		}`, false, []Rule{{
			NameScope:  ".",
			OutputFile: DefaultLogFilename,
			Format:     CommonLogFormat,
			Class:      response.Denial,
		}}},
	}
	for i, test := range tests {
		c := caddy.NewTestController("dns", test.inputLogRules)
		actualLogRules, err := logParse(c)

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}
		if len(actualLogRules) != len(test.expectedLogRules) {
			t.Fatalf("Test %d expected %d no of Log rules, but got %d ",
				i, len(test.expectedLogRules), len(actualLogRules))
		}
		for j, actualLogRule := range actualLogRules {

			if actualLogRule.NameScope != test.expectedLogRules[j].NameScope {
				t.Errorf("Test %d expected %dth LogRule NameScope to be  %s  , but got %s",
					i, j, test.expectedLogRules[j].NameScope, actualLogRule.NameScope)
			}

			if actualLogRule.OutputFile != test.expectedLogRules[j].OutputFile {
				t.Errorf("Test %d expected %dth LogRule OutputFile to be  %s  , but got %s",
					i, j, test.expectedLogRules[j].OutputFile, actualLogRule.OutputFile)
			}

			if actualLogRule.Format != test.expectedLogRules[j].Format {
				t.Errorf("Test %d expected %dth LogRule Format to be  %s  , but got %s",
					i, j, test.expectedLogRules[j].Format, actualLogRule.Format)
			}

			if actualLogRule.Class != test.expectedLogRules[j].Class {
				t.Errorf("Test %d expected %dth LogRule Class to be  %s  , but got %s",
					i, j, test.expectedLogRules[j].Class, actualLogRule.Class)
			}

			if actualLogRule.MaxAge != test.expectedLogRules[j].MaxAge {
				t.Errorf("Test %d expected %dth LogRule MaxAge to be  %d  , but got %d",
					i, j, test.expectedLogRules[j].MaxAge, actualLogRule.MaxAge)
			}

			if actualLogRule.MaxBackups != test.expectedLogRules[j].MaxBackups {
				t.Errorf("Test %d expected %dth LogRule MaxBackups to be  %d  , but got %d",
					i, j, test.expectedLogRules[j].MaxBackups, actualLogRule.MaxBackups)
			}

			if actualLogRule.MaxSize != test.expectedLogRules[j].MaxSize {
				t.Errorf("Test %d expected %dth LogRule MaxSize to be  %d  , but got %d",
					i, j, test.expectedLogRules[j].MaxSize, actualLogRule.MaxSize)
			}
		}
	}

}
