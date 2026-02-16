package hba

import "strings"

// Rule represents one entry in pg_hba.conf.
// Order matters: PostgreSQL uses the first matching rule.
type Rule struct {
	Type     string // local, host, hostssl, hostnossl, hostgssenc, hostnogssenc
	Database string // database name or "all", "sameuser", "samerole", "replication"
	User     string // user name or "all"
	Address  string // IP/CIDR or "samehost", "samenet"; "-" for local
	Netmask  string // optional: legacy IP netmask (e.g. 255.255.255.0); empty when using CIDR
	Method   string // trust, reject, scram-sha-256, md5, etc.; may include auth options
}

// MatchesUser returns true if the rule matches the given user and optional database.
// If db is empty, only user is matched. User and db are compared case-sensitively.
func (r Rule) MatchesUser(user, db string) bool {
	if user == "" {
		return false
	}
	if r.User != user {
		return false
	}
	if db != "" && r.Database != db {
		return false
	}
	return true
}

// MatchesAddress returns true if the rule's address matches the given addr.
// Matches exact Address, or addr as CIDR (e.g. 10.0.1.7 matches 10.0.1.7 or 10.0.1.7/32).
func (r Rule) MatchesAddress(addr string) bool {
	if addr == "" {
		return false
	}
	if r.Address == addr {
		return true
	}
	// Allow 10.0.1.7 to match 10.0.1.7/32
	if !strings.Contains(addr, "/") && r.Address == addr+"/32" {
		return true
	}
	if strings.HasSuffix(addr, "/32") && r.Address == strings.TrimSuffix(addr, "/32") {
		return true
	}
	return false
}
