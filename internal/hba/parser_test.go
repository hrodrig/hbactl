package hba

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := `# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
# comment line
host    mydb            app             192.168.1.0/24           md5
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rules, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if want := 4; len(rules) != want {
		t.Fatalf("got %d rules, want %d", len(rules), want)
	}
	checkRule0(t, rules[0])
	checkRule1(t, rules[1])
	checkRule2(t, rules[2])
	checkRule3(t, rules[3])
}

func checkRule0(t *testing.T, r Rule) {
	t.Helper()
	if r.Type != "local" || r.Database != "all" || r.User != "all" || r.Address != "-" || r.Method != "trust" {
		t.Errorf("rule 0: got %+v", r)
	}
}

func checkRule1(t *testing.T, r Rule) {
	t.Helper()
	if r.Type != "host" || r.Address != "127.0.0.1/32" || r.Method != "scram-sha-256" {
		t.Errorf("rule 1: got %+v", r)
	}
}

func checkRule2(t *testing.T, r Rule) {
	t.Helper()
	if r.Address != "::1/128" {
		t.Errorf("rule 2 address: got %q", r.Address)
	}
}

func checkRule3(t *testing.T, r Rule) {
	t.Helper()
	if r.Database != "mydb" || r.User != "app" || r.Method != "md5" {
		t.Errorf("rule 3: got %+v", r)
	}
}

func TestParseFile_emptyAndComments(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := `
# only comments
# TYPE DATABASE USER METHOD

local all all trust
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rules, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
	if rules[0].Type != "local" || rules[0].Method != "trust" {
		t.Errorf("got %+v", rules[0])
	}
}

func TestParseFile_legacyNetmask(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := "host\tall\tall\t127.0.0.1\t255.255.255.255\ttrust\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rules, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
	if rules[0].Address != "127.0.0.1" || rules[0].Netmask != "255.255.255.255" || rules[0].Method != "trust" {
		t.Errorf("legacy netmask: got %+v", rules[0])
	}
}

func TestParseFile_hostssl(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := "hostssl\tall\tall\t0.0.0.0/0\tscram-sha-256\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rules, err := ParseFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 {
		t.Fatalf("got %d rules, want 1", len(rules))
	}
	if rules[0].Type != "hostssl" || rules[0].Address != "0.0.0.0/0" {
		t.Errorf("got %+v", rules[0])
	}
}
