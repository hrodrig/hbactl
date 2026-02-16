package hba

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := "local\tall\tall\ttrust\nhost\tall\tall\t127.0.0.1/32\tscram-sha-256\nhost\tmydb\tapp\t192.168.1.0/24\tmd5\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	// Remove middle line (second rule)
	if err := RemoveLine(path, 2); err != nil {
		t.Fatalf("RemoveLine(2): %v", err)
	}
	data, _ := os.ReadFile(path)
	got := string(data)
	want := "local\tall\tall\ttrust\nhost\tmydb\tapp\t192.168.1.0/24\tmd5\n"
	if got != want {
		t.Errorf("after remove line 2: got %q, want %q", got, want)
	}

	// Remove first line
	if err := RemoveLine(path, 1); err != nil {
		t.Fatalf("RemoveLine(1): %v", err)
	}
	data, _ = os.ReadFile(path)
	got = string(data)
	want = "host\tmydb\tapp\t192.168.1.0/24\tmd5\n"
	if got != want {
		t.Errorf("after remove line 1: got %q, want %q", got, want)
	}

	// Remove last remaining line
	if err := RemoveLine(path, 1); err != nil {
		t.Fatalf("RemoveLine(1): %v", err)
	}
	data, _ = os.ReadFile(path)
	got = string(data)
	want = "\n"
	if got != want {
		t.Errorf("after remove last line: got %q, want %q", got, want)
	}
}

func TestRemoveLine_outOfRange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	// No trailing newline so split gives exactly one line
	if err := os.WriteFile(path, []byte("local\tall\tall\ttrust"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := RemoveLine(path, 0); err == nil {
		t.Error("RemoveLine(0) should fail")
	}
	if err := RemoveLine(path, 2); err == nil {
		t.Error("RemoveLine(2) should fail for single-line file")
	}
}

func TestParseFileWithLineNumbers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pg_hba.conf")
	content := "# comment\nlocal\tall\tall\ttrust\nhost\tall\tall\t127.0.0.1/32\tscram-sha-256\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	rwl, err := ParseFileWithLineNumbers(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rwl) != 2 {
		t.Fatalf("got %d rules, want 2", len(rwl))
	}
	if rwl[0].LineNo != 2 || rwl[0].Index != 1 {
		t.Errorf("first rule: LineNo=%d Index=%d, want LineNo=2 Index=1", rwl[0].LineNo, rwl[0].Index)
	}
	if rwl[1].LineNo != 3 || rwl[1].Index != 2 {
		t.Errorf("second rule: LineNo=%d Index=%d, want LineNo=3 Index=2", rwl[1].LineNo, rwl[1].Index)
	}
}
