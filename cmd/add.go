package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hrodrig/hbactl/internal/hba"
	"github.com/hrodrig/hbactl/internal/pg"
	"github.com/spf13/cobra"
)

var (
	addType     string
	addDB       string
	addUser     string
	addAddr     string
	addNetmask  string
	addMethod   string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a rule to pg_hba.conf",
	Long:  "Appends a new rule to pg_hba.conf. Creates a backup before writing. Run 'hbactl reload' to apply changes.",
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addType, "type", "", "Rule type: local, host, hostssl, hostnossl, hostgssenc, hostnogssenc")
	addCmd.Flags().StringVar(&addDB, "db", "all", "Database (e.g. all, sameuser, or name)")
	addCmd.Flags().StringVar(&addUser, "user", "all", "User (e.g. all or name)")
	addCmd.Flags().StringVar(&addAddr, "addr", "", "Address: CIDR (e.g. 127.0.0.1/32), samehost, samenet; use - for local")
	addCmd.Flags().StringVar(&addNetmask, "netmask", "", "Optional: legacy netmask (e.g. 255.255.255.0)")
	addCmd.Flags().StringVar(&addMethod, "method", "", "Auth method: trust, reject, scram-sha-256, md5, ident, etc.")
	_ = addCmd.MarkFlagRequired("type")
	_ = addCmd.MarkFlagRequired("method")
}

func runAdd(cmd *cobra.Command, _ []string) error {
	typ := strings.ToLower(strings.TrimSpace(addType))
	method := strings.TrimSpace(addMethod)
	db := strings.TrimSpace(addDB)
	if db == "" {
		db = "all"
	}
	user := strings.TrimSpace(addUser)
	if user == "" {
		user = "all"
	}
	addr := strings.TrimSpace(addAddr)
	netmask := strings.TrimSpace(addNetmask)

	if !hba.LocalType(typ) && !hba.HostType(typ) {
		return fmt.Errorf("invalid type %q; use one of: local, host, hostssl, hostnossl, hostgssenc, hostnogssenc", typ)
	}
	if hba.LocalType(typ) {
		addr = "-"
		netmask = ""
	} else {
		if addr == "" {
			return fmt.Errorf("addr is required for type %s (e.g. 127.0.0.1/32, samehost)", typ)
		}
	}

	path := filePath()
	if path == "" {
		conn := connString()
		if conn == "" {
			return fmt.Errorf("no connection: set DATABASE_URL or use --conn (or pass path with --file)")
		}
		ctx := context.Background()
		client, err := pg.NewClient(ctx, conn)
		if err != nil {
			return fmt.Errorf("could not connect to PostgreSQL: %w", err)
		}
		defer client.Close()
		var p string
		p, err = client.HBAFilePath(ctx)
		if err != nil {
			return fmt.Errorf("could not locate pg_hba.conf. Is PostgreSQL running? %w", err)
		}
		path = p
	}

	backupPath, err := hba.Backup(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("backup failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Backup created at: %s\n", backupPath)

	rule := hba.Rule{Type: typ, Database: db, User: user, Address: addr, Netmask: netmask, Method: method}
	if err := hba.AppendRule(path, rule); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("failed to append rule: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Success: New rule added to %s. Run 'hbactl reload' to apply changes.\n", path)
	return nil
}
