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
	removeIndex  int
	removeUser   string
	removeDB     string
	removeAddr   string
	removeDryRun bool
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove rule(s) from pg_hba.conf by index or by criteria",
	Long:  "Removes one rule by --index, or all matching rules by --user (optional --db) or --addr. Creates a backup before editing. Use --dry-run to preview. Run 'hbactl reload' after to apply changes.",
	RunE:  runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().IntVar(&removeIndex, "index", 0, "1-based rule index to remove (see first column of 'hbactl list')")
	removeCmd.Flags().StringVar(&removeUser, "user", "", "Remove all rules for this user; combine with --db to limit by database")
	removeCmd.Flags().StringVar(&removeDB, "db", "", "When used with --user, only remove rules for this database")
	removeCmd.Flags().StringVar(&removeAddr, "addr", "", "Remove all rules matching this address (e.g. 10.0.1.7 or 10.0.1.7/32)")
	removeCmd.Flags().BoolVar(&removeDryRun, "dry-run", false, "Print the rule(s) that would be removed without writing or creating backup")
}

func runRemove(cmd *cobra.Command, _ []string) error {
	byIndex := removeIndex >= 1
	byUser := strings.TrimSpace(removeUser) != ""
	byAddr := strings.TrimSpace(removeAddr) != ""

	if byIndex && (byUser || byAddr) {
		return fmt.Errorf("use either --index or criteria (--user/--addr), not both")
	}
	if byUser && byAddr {
		return fmt.Errorf("use only one of --user or --addr per run")
	}
	if !byIndex && !byUser && !byAddr {
		return fmt.Errorf("specify --index N, or --user <name> [--db <name>], or --addr <address>")
	}
	if byIndex && removeIndex < 1 {
		return fmt.Errorf("--index must be >= 1")
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
		p, err := client.HBAFilePath(ctx)
		if err != nil {
			return fmt.Errorf("could not locate pg_hba.conf. Is PostgreSQL running? %w", err)
		}
		path = p
	}

	rwl, err := hba.ParseFileWithLineNumbers(path)
	if err != nil {
		return fmt.Errorf("could not read file (try running with sudo?): %w", err)
	}

	var toRemove []hba.RuleWithLine
	if byIndex {
		for i := range rwl {
			if rwl[i].Index == removeIndex {
				toRemove = append(toRemove, rwl[i])
				break
			}
		}
		if len(toRemove) == 0 {
			return fmt.Errorf("no rule at index %d (file has %d rule(s)); run 'hbactl list' to see indices", removeIndex, len(rwl))
		}
	} else {
		user := strings.TrimSpace(removeUser)
		db := strings.TrimSpace(removeDB)
		addr := strings.TrimSpace(removeAddr)
		for i := range rwl {
			r := rwl[i].Rule
			matches := false
			if byUser && r.MatchesUser(user, db) {
				matches = true
			}
			if byAddr && r.MatchesAddress(addr) {
				matches = true
			}
			if matches {
				toRemove = append(toRemove, rwl[i])
			}
		}
		if len(toRemove) == 0 {
			crit := ""
			if byUser {
				crit = fmt.Sprintf("user %q", user)
				if db != "" {
					crit += fmt.Sprintf(" db %q", db)
				}
			} else {
				crit = fmt.Sprintf("addr %q", addr)
			}
			return fmt.Errorf("no rules matching %s; run 'hbactl list' to inspect", crit)
		}
	}

	if removeDryRun {
		fmt.Fprintf(os.Stdout, "dry-run: would remove %d rule(s) from %s:\n", len(toRemove), path)
		for _, x := range toRemove {
			fmt.Fprintf(os.Stdout, "  #%d (line %d): %s\n", x.Index, x.LineNo, x.Rule.Line())
		}
		return nil
	}

	backupPath, err := hba.Backup(path)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("backup failed: %w", err)
	}
	fmt.Fprintf(os.Stderr, "Backup created at: %s\n", backupPath)

	lineNos := make([]int, len(toRemove))
	for i := range toRemove {
		lineNos[i] = toRemove[i].LineNo
	}
	if err := hba.RemoveLines(path, lineNos); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("insufficient permissions to write to pg_hba.conf. Try running with sudo")
		}
		return fmt.Errorf("remove failed: %w", err)
	}

	if len(toRemove) == 1 {
		fmt.Fprintf(os.Stdout, "Success: Rule #%d removed from %s. Run 'hbactl reload' to apply changes.\n", toRemove[0].Index, path)
	} else {
		fmt.Fprintf(os.Stdout, "Success: %d rule(s) removed from %s. Run 'hbactl reload' to apply changes.\n", len(toRemove), path)
	}
	return nil
}
