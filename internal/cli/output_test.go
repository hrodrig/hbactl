package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hrodrig/hbactl/internal/hba"
)

func TestWriteRulesTable(t *testing.T) {
	rules := []hba.Rule{
		{Type: "local", Database: "all", User: "all", Address: "-", Method: "trust"},
		{Type: "host", Database: "all", User: "all", Address: "127.0.0.1/32", Method: "scram-sha-256"},
	}
	var buf bytes.Buffer
	WriteRulesTable(&buf, rules)
	out := buf.String()
	if !strings.Contains(out, "TYPE") || !strings.Contains(out, "DATABASE") {
		t.Errorf("table missing header: %s", out)
	}
	if !strings.Contains(out, "local") || !strings.Contains(out, "trust") {
		t.Errorf("table missing rule data: %s", out)
	}
	if !strings.Contains(out, "127.0.0.1/32") {
		t.Errorf("table missing address: %s", out)
	}
}

func TestWriteRulesTable_empty(t *testing.T) {
	var buf bytes.Buffer
	WriteRulesTable(&buf, nil)
	out := buf.String()
	if !strings.Contains(out, "TYPE") {
		t.Errorf("empty table should still have header: %s", out)
	}
}
