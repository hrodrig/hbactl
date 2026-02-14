package hba

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
