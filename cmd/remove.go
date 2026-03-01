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

type removeMode struct {
	byIndex bool
	byUser  bool
	byAddr  bool
	user    string
	db      string
	addr    string
}

func validateRemoveFlags() (removeMode, error) {
	byIndex := removeIndex >= 1
	byUser := strings.TrimSpace(removeUser) != ""
	byAddr := strings.TrimSpace(removeAddr) != ""

	if byIndex && (byUser || byAddr) {
		return removeMode{}, fmt.Errorf("use either --index or criteria (--user/--addr), not both")
	}
	if byUser && byAddr {
		return removeMode{}, fmt.Errorf("use only one of --user or --addr per run")
	}
	if !byIndex && !byUser && !byAddr {
		return removeMode{}, fmt.Errorf("specify --index N, or --user <name> [--db <name>], or --addr <address>")
	}
	if byIndex && removeIndex < 1 {
		return removeMode{}, fmt.Errorf("--index must be >= 1")
	}
	return removeMode{
		byIndex: byIndex,
		byUser:  byUser,
		byAddr:  byAddr,
		user:    strings.TrimSpace(removeUser),
		db:      strings.TrimSpace(removeDB),
		addr:    strings.TrimSpace(removeAddr),
	}, nil
}

func resolveRemovePath() (string, error) {
	path := filePath()
	if path != "" {
		return path, nil
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

func findRulesToRemove(rwl []hba.RuleWithLine, mode removeMode) ([]hba.RuleWithLine, error) {
	if mode.byIndex {
		for i := range rwl {
			if rwl[i].Index == removeIndex {
				return []hba.RuleWithLine{rwl[i]}, nil
			}
		}
		return nil, fmt.Errorf("no rule at index %d (file has %d rule(s)); run 'hbactl list' to see indices", removeIndex, len(rwl))
	}
	var out []hba.RuleWithLine
	for i := range rwl {
		r := rwl[i].Rule
		matches := (mode.byUser && r.MatchesUser(mode.user, mode.db)) || (mode.byAddr && r.MatchesAddress(mode.addr))
		if matches {
			out = append(out, rwl[i])
		}
	}
	if len(out) == 0 {
		crit := removeCriteriaString(mode)
		return nil, fmt.Errorf("no rules matching %s; run 'hbactl list' to inspect", crit)
	}
	return out, nil
}

func removeCriteriaString(mode removeMode) string {
	if mode.byUser {
		s := fmt.Sprintf("user %q", mode.user)
		if mode.db != "" {
			s += fmt.Sprintf(" db %q", mode.db)
		}
		return s
	}
	return fmt.Sprintf("addr %q", mode.addr)
}

func runRemove(cmd *cobra.Command, _ []string) error {
	mode, err := validateRemoveFlags()
	if err != nil {
		return err
	}
	path, err := resolveRemovePath()
	if err != nil {
		return err
	}
	rwl, err := hba.ParseFileWithLineNumbers(path)
	if err != nil {
		return fmt.Errorf("could not read file (try running with sudo?): %w", err)
	}
	toRemove, err := findRulesToRemove(rwl, mode)
	if err != nil {
		return err
	}
	if removeDryRun {
		printRemoveDryRun(toRemove, path)
		return nil
	}
	if err := doRemove(path, toRemove); err != nil {
		return err
	}
	printRemoveSuccess(toRemove, path)
	return nil
}

func printRemoveDryRun(toRemove []hba.RuleWithLine, path string) {
	fmt.Fprintf(os.Stdout, "dry-run: would remove %d rule(s) from %s:\n", len(toRemove), path)
	for _, x := range toRemove {
		fmt.Fprintf(os.Stdout, "  #%d (line %d): %s\n", x.Index, x.LineNo, x.Rule.Line())
	}
}

func doRemove(path string, toRemove []hba.RuleWithLine) error {
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
	return nil
}

func printRemoveSuccess(toRemove []hba.RuleWithLine, path string) {
	if len(toRemove) == 1 {
		fmt.Fprintf(os.Stdout, "Success: Rule #%d removed from %s. Run 'hbactl reload' to apply changes.\n", toRemove[0].Index, path)
	} else {
		fmt.Fprintf(os.Stdout, "Success: %d rule(s) removed from %s. Run 'hbactl reload' to apply changes.\n", len(toRemove), path)
	}
}
