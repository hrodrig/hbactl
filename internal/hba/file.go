package hba

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// InsertRuleAfterUser inserts the rule after the last rule whose User equals afterUser.
// If afterUser is empty or no such rule exists, the rule is appended at the end.
// Call Backup before this if you want a backup.
func InsertRuleAfterUser(path string, r Rule, afterUser string) error {
	line := r.Line()
	if line == "" {
		return fmt.Errorf("invalid rule type %q", r.Type)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := string(data)
	lines := strings.Split(content, "\n")
	// If file ends with newline, lines has a trailing ""; we'll preserve that when inserting.
	lastIdx := -1
	for i := range lines {
		ln := strings.TrimSpace(lines[i])
		if j := strings.Index(ln, "#"); j >= 0 {
			ln = strings.TrimSpace(ln[:j])
		}
		if ln == "" {
			continue
		}
		parsed, ok := parseLine(ln)
		if !ok {
			continue
		}
		if parsed.User == afterUser {
			lastIdx = i
		}
	}
	var newLines []string
	if lastIdx < 0 {
		// Append: drop trailing empty from split so we don't add an extra blank line
		trimmed := lines
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			trimmed = lines[:len(lines)-1]
		}
		newLines = append(trimmed, line)
	} else {
		newLines = append(lines[:lastIdx+1], append([]string{line}, lines[lastIdx+1:]...)...)
	}
	out := strings.Join(newLines, "\n")
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0644)
}

// Backup copies path to path.bak (or path.bak.<timestamp> if path.bak exists). Returns the backup path.
func Backup(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	backupPath := path + ".bak"
	if _, err := os.Stat(backupPath); err == nil {
		backupPath = fmt.Sprintf("%s.bak.%s", path, time.Now().Format("20060102-150405"))
	}
	return backupPath, os.WriteFile(backupPath, data, 0644)
}

// Line returns the pg_hba.conf line for the rule (one line, no newline).
func (r Rule) Line() string {
	if localTypes[strings.ToLower(r.Type)] {
		return fmt.Sprintf("%s\t%s\t%s\t%s", r.Type, r.Database, r.User, r.Method)
	}
	if hostTypes[strings.ToLower(r.Type)] {
		if r.Netmask != "" {
			return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", r.Type, r.Database, r.User, r.Address, r.Netmask, r.Method)
		}
		return fmt.Sprintf("%s\t%s\t%s\t%s\t%s", r.Type, r.Database, r.User, r.Address, r.Method)
	}
	return ""
}

// AppendRule appends the rule line to the file, prefixed with newline if file is not empty.
func AppendRule(path string, r Rule) error {
	line := r.Line()
	if line == "" {
		return fmt.Errorf("invalid rule type %q", r.Type)
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if info.Size() > 0 {
		if _, err := f.WriteString("\n"); err != nil {
			return err
		}
	}
	_, err = f.WriteString(line + "\n")
	return err
}
