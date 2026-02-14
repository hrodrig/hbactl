# hbactl

[![version](https://img.shields.io/badge/version-0.1.0-blue)](https://github.com/hrodrig/hbactl/releases) [![license](https://img.shields.io/badge/license-MIT-green)](LICENSE)

`hbactl` is a lightweight CLI tool written in **Go** designed to manage your PostgreSQL Host-Based Authentication (`pg_hba.conf`) file safely and efficiently.

## Features

- **Auto-Discovery**: Locates `pg_hba.conf` via the running Postgres instance, or use `--file` to pass the path.
- **Safety First**: Backup before every edit; validate syntax with `hbactl check` (uses `pg_hba_file_rules`).
- **Reload**: Apply changes with `hbactl reload` (`pg_reload_conf()`), no restart.
- **Single Binary**: One executable; no runtime dependencies.
- **Formats**: Supports both CIDR (e.g. `192.168.1.0/24`) and legacy IP+netmask in `list` and `add`.

## Installation

**From Go (recommended):**

```bash
go install github.com/hrodrig/hbactl@latest
```

**From GitHub Releases:**  
Pre-built binaries for Linux, macOS (darwin), and Windows (amd64, arm64) are published on [Releases](https://github.com/hrodrig/hbactl/releases). Download the archive for your OS/arch and unpack the `hbactl` binary.

**Build from source:**

```bash
git clone https://github.com/hrodrig/hbactl.git
cd hbactl
go build -o hbactl .
```

## Usage

Global flags (optional):

- **`-c` / `--conn`** — PostgreSQL connection string (default: `DATABASE_URL`).
- **`-f` / `--file`** — Path to `pg_hba.conf` (for `list` and `add`; can avoid connection for `list`).

### List current rules

Displays a formatted table of your rules. Supports **`--sort`** by column: `type`, `database`, `user`, `address`, `method` (display only; file order is unchanged).

```bash
hbactl list
hbactl list -f /path/to/pg_hba.conf              # no connection needed
hbactl list --sort type
```

### Add a new rule

Creates a **backup** (`.bak` or `.bak.<timestamp>`) then appends the rule. Requires connection or `--file`.

```bash
hbactl add --type host --db all --user app --addr 192.168.1.100/32 --method scram-sha-256
hbactl add --type host --db all --user all --addr 10.0.0.0/24 --netmask 255.255.255.0 --method md5   # legacy format
hbactl add -f /path/to/pg_hba.conf --type host --db all --user all --addr 127.0.0.1/32 --method trust
```

Flags: **`--type`** (required: `local`, `host`, `hostssl`, …), **`--db`**, **`--user`**, **`--addr`** (required for host types), **`--netmask`** (optional, legacy), **`--method`** (required).

### Check for errors

Uses `pg_hba_file_rules` to report syntax errors. Requires a connection.

```bash
hbactl check
```

### Reload configuration

Runs `pg_reload_conf()` so the server picks up changes without restart. Requires a connection.

```bash
hbactl reload
```

## Connection

By default, `hbactl` connects to PostgreSQL using the `DATABASE_URL` environment variable or the `--conn` / `-c` flag. Ensure your user has permission to read the HBA file and run `pg_reload_conf()`.

### Multiple servers / explicit file path

Use `--file` / `-f` to pass the path to `pg_hba.conf` explicitly. With `--file`, `list` can run **without a connection** (it only reads the file). Handy for:

- Managing several servers: pass a different path (or `--conn` + path) per run.
- Inspecting a local copy of the file (e.g. from another host).

```bash
# List rules from a file (no DB connection)
hbactl list -f /opt/homebrew/var/postgresql@16/pg_hba.conf

# Per-server: connect to server and optionally override path
hbactl list --conn "postgres://user@server1:5432/postgres" -f /var/lib/pgsql/16/data/pg_hba.conf
```

### Checking if PostgreSQL is running

- **Port**: Something should be listening on `5432` (e.g. `lsof -i :5432`, `ss -tlnp | grep 5432`, or `nc -z localhost 5432`).
- **Client**: If you have `psql` or `pg_isready` installed, run `pg_isready -h localhost -p 5432`.
- **macOS (Homebrew)**: `brew services list | grep postgres`.
- **Linux (systemd)**: `systemctl status postgresql` or `systemctl status postgresql@16` (depending on distro and version).

If nothing is listening on 5432 and no Postgres service is listed, you don’t have a running server.

### Running PostgreSQL for testing

**Linux (Debian/Ubuntu):**

```bash
sudo apt update && sudo apt install -y postgresql postgresql-client
sudo systemctl start postgresql
sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'postgres';"
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres"
hbactl list
```

**Linux (Fedora/RHEL):**

```bash
sudo dnf install -y postgresql-server postgresql
sudo postgresql-setup --initdb
sudo systemctl start postgresql
export DATABASE_URL="postgres://postgres@localhost:5432/postgres"   # peer auth by default
hbactl list
```

**Homebrew (macOS):**

```bash
brew install postgresql@16
brew services start postgresql@16
export DATABASE_URL="postgres://$(whoami)@localhost:5432/postgres"
hbactl list
```

**Docker:**

```bash
docker run -d --name pg -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/postgres"
hbactl list
```

**Unit tests** (parser and table output) do not require a running PostgreSQL server: run `go test ./...`.

## Releasing

Releases and tags are made **only from the `main` branch**. Work on `develop` (or feature branches), merge to `main` when ready, then from `main`:

```bash
git checkout main
git pull
git tag v0.1.0
make release
git push origin v0.1.0
```

`make release` will fail if you are not on `main`. Use `make snapshot` on any branch to test builds locally.

## License

MIT
