package hba

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInsertRuleAfterUser(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := "host\tall\tpepe\t192.168.1.1/32\tmd5\nhost\tall\tjuan\t10.0.0.1/32\ttrust\nhost\tall\tpepe\t192.168.1.2/32\tscram-sha-256\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Insert new pepe rule after last pepe -> should go after 192.168.1.2/32 line
	rule := Rule{Type: "host", Database: "all", User: "pepe", Address: "10.0.0.5/32", Method: "md5"}
	if err := InsertRuleAfterUser(path, rule, "pepe"); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	got := string(data)
	want := "host\tall\tpepe\t192.168.1.1/32\tmd5\nhost\tall\tjuan\t10.0.0.1/32\ttrust\nhost\tall\tpepe\t192.168.1.2/32\tscram-sha-256\nhost\tall\tpepe\t10.0.0.5/32\tmd5\n"
	if got != want {
		t.Errorf("insert after pepe:\ngot:\n%q\nwant:\n%q", got, want)
	}

	// Insert with after-user that has no rule -> append at end
	rule2 := Rule{Type: "host", Database: "all", User: "newuser", Address: "127.0.0.1/32", Method: "trust"}
	if err := InsertRuleAfterUser(path, rule2, "nonexistent"); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	got = string(data)
	if len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("insert after nonexistent should append; got prefix:\n%q", got[:len(want)])
	}
	if got[len(want):] != "host\tall\tnewuser\t127.0.0.1/32\ttrust\n" {
		t.Errorf("appended line: got %q", got[len(want):])
	}
}
