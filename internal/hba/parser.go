package hba

import (
	"bufio"
	"os"
	"strings"
)

// localTypes are connection types that have no address field (4 fields: type, database, user, method).
var localTypes = map[string]bool{
	"local": true,
}

// hostTypes are connection types that have an address field (5+ fields: type, database, user, address, method [, options]).
var hostTypes = map[string]bool{
	"host":        true,
	"hostssl":     true,
	"hostnossl":   true,
	"hostgssenc":  true,
	"hostnogssenc": true,
}

// LocalType returns true if typ is a local connection type (no address).
func LocalType(typ string) bool { return localTypes[typ] }

// HostType returns true if typ is a host connection type (has address).
func HostType(typ string) bool { return hostTypes[typ] }

// ParseFile reads path and returns parsed rules. Comment and empty lines are skipped.
func ParseFile(path string) ([]Rule, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var rules []Rule
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Strip inline comment
		if i := strings.Index(line, "#"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}
		rule, ok := parseLine(line)
		if !ok {
			continue // skip malformed or non-rule lines
		}
		rules = append(rules, rule)
	}
	return rules, scanner.Err()
}

// parseLine parses one line into a Rule. Returns ok=false if the line is not a valid rule.
func parseLine(line string) (Rule, bool) {
	fields := splitFields(line)
	if len(fields) < 4 {
		return Rule{}, false
	}
	typ := strings.ToLower(fields[0])
	if localTypes[typ] {
		// local database user method [options]
		if len(fields) < 4 {
			return Rule{}, false
		}
		method := fields[3]
		if len(fields) > 4 {
			method = strings.Join(fields[3:], " ")
		}
		return Rule{
			Type:     typ,
			Database: fields[1],
			User:     fields[2],
			Address:  "-",
			Method:   method,
		}, true
	}
	if hostTypes[typ] {
		// host database user address method [options]
		// or legacy: host database user IP-ADDRESS IP-MASK method [options] (6+ fields)
		if len(fields) < 5 {
			return Rule{}, false
		}
		addr := fields[3]
		netmask := ""
		methodStart := 4
		if len(fields) >= 6 && looksLikeNetmask(fields[4]) {
			// Legacy format: IP + netmask
			netmask = fields[4]
			methodStart = 5
		}
		method := fields[methodStart]
		if len(fields) > methodStart+1 {
			method = strings.Join(fields[methodStart:], " ")
		}
		return Rule{
			Type:     typ,
			Database: fields[1],
			User:     fields[2],
			Address:  addr,
			Netmask:  netmask,
			Method:   method,
		}, true
	}
	return Rule{}, false
}

// looksLikeNetmask returns true if s looks like a legacy netmask (dotted decimal or IPv6 hex).
func looksLikeNetmask(s string) bool {
	if s == "" {
		return false
	}
	// IPv4: 255.255.255.0 or 255.0.0.0
	if strings.Contains(s, ".") {
		return true
	}
	// IPv6: ffff:ffff:...
	if strings.Contains(s, ":") {
		return true
	}
	return false
}

// splitFields splits by whitespace, respecting double-quoted segments (one field may contain spaces).
func splitFields(s string) []string {
	var fields []string
	var buf strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			inQuote = !inQuote
		case inQuote:
			buf.WriteByte(c)
		case c == ' ' || c == '\t':
			if buf.Len() > 0 {
				fields = append(fields, buf.String())
				buf.Reset()
			}
		default:
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		fields = append(fields, buf.String())
	}
	return fields
}
