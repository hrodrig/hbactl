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
	addType      string
	addDB        string
	addUser      string
	addAddr      string
	addNetmask   string
	addMethod    string
	addIdentMap  string
	addDryRun    bool
	addAfterUser string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a rule to pg_hba.conf",
	Long:  "Appends a new rule to pg_hba.conf. Creates a backup before writing. Run 'hbactl reload' to apply changes. Use --dry-run to preview without writing.",
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addType, "type", "", "Rule type: local, host, hostssl, hostnossl, hostgssenc, hostnogssenc")
	addCmd.Flags().StringVar(&addDB, "db", "all", "Database (e.g. all, sameuser, or name)")
	addCmd.Flags().StringVar(&addUser, "user", "all", "User: all, a user name, or a group/role with + prefix (e.g. +admins)")
	addCmd.Flags().StringVar(&addAddr, "addr", "", "Address: CIDR (e.g. 127.0.0.1/32), samehost, samenet; use - for local")
	addCmd.Flags().StringVar(&addNetmask, "netmask", "", "Optional: legacy netmask (e.g. 255.255.255.0)")
	addCmd.Flags().StringVar(&addMethod, "method", "", "Auth method: trust, reject, scram-sha-256, md5, ident, etc.")
	addCmd.Flags().StringVar(&addIdentMap, "ident-map", "", "For method ident: username map name (e.g. my_ident_map → writes 'ident my_ident_map')")
	addCmd.Flags().BoolVar(&addDryRun, "dry-run", false, "Print the line that would be added without writing or creating backup")
	addCmd.Flags().StringVar(&addAfterUser, "after-user", "", "Insert after the last rule for this user (keeps rules grouped by user); default appends at end")
	_ = addCmd.MarkFlagRequired("type")
	_ = addCmd.MarkFlagRequired("method")
}

func normalizeAddFlags() (hba.Rule, error) {
	typ := strings.ToLower(strings.TrimSpace(addType))
	method := strings.TrimSpace(addMethod)
	identMap := strings.TrimSpace(addIdentMap)
	if identMap != "" && strings.ToLower(method) == "ident" {
		method = "ident " + identMap
	}
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
		return hba.Rule{}, fmt.Errorf("invalid type %q; use one of: local, host, hostssl, hostnossl, hostgssenc, hostnogssenc", typ)
	}
	if hba.LocalType(typ) {
		addr = "-"
		netmask = ""
	} else if addr == "" {
		return hba.Rule{}, fmt.Errorf("addr is required for type %s (e.g. 127.0.0.1/32, samehost)", typ)
	}
	return hba.Rule{Type: typ, Database: db, User: user, Address: addr, Netmask: netmask, Method: method}, nil
}

func resolveAddPath(requirePath bool) (string, error) {
	path := filePath()
	if path != "" {
		return path, nil
	}
	if !requirePath {
		return "(path from --file or connection)", nil
	}
	conn := connString()
	if conn == "" {
		return "", fmt.Errorf("no connection: set DATABASE_URL or use --conn (or pass path with --file)")
	}
	ctx := context.Background()
	client, err := pg.NewClient(ctx, conn)
	if err != nil {
		return "", fmt.Errorf("could not connect to PostgreSQL: %w", err)
	}
	defer client.Close()
	p, err := client.HBAFilePath(ctx)
	if err != nil {
		return "", fmt.Errorf("could not locate pg_hba.conf. Is PostgreSQL running? %w", err)
	}
	return p, nil
}

func printAddDryRun(rule hba.Rule, path, afterUser string) error {
	if rule.Line() == "" {
		return fmt.Errorf("invalid rule type %q", rule.Type)
	}
	if afterUser != "" {
		fmt.Fprintf(os.Stdout, "dry-run: would insert after last rule for user %q in %s:\n%s\n", afterUser, path, rule.Line())
	} else {
		fmt.Fprintf(os.Stdout, "dry-run: would append to %s:\n%s\n", path, rule.Line())
	}
	return nil
}

func doAdd(path string, rule hba.Rule, afterUser string) error {
	backupPath, err := hba.Backup(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("backup failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Backup created at: %s\n", backupPath)
	if afterUser != "" {
		if err := hba.InsertRuleAfterUser(path, rule, afterUser); err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
			}
			return fmt.Errorf("failed to insert rule after user %q: %w", afterUser, err)
		}
		return nil
	}
	if err := hba.AppendRule(path, rule); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("failed to append rule: %w", err)
	}
	return nil
}

func runAdd(cmd *cobra.Command, _ []string) error {
	rule, err := normalizeAddFlags()
	if err != nil {
		return err
	}
	path, err := resolveAddPath(!addDryRun)
	if err != nil {
		return err
	}
	if addDryRun {
		return printAddDryRun(rule, path, strings.TrimSpace(addAfterUser))
	}
	if err := doAdd(path, rule, strings.TrimSpace(addAfterUser)); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Success: New rule added to %s. Run 'hbactl reload' to apply changes.\n", path)
	return nil
}
