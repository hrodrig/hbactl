package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hrodrig/hbactl/internal/cli"
	"github.com/hrodrig/hbactl/internal/hba"
	"github.com/hrodrig/hbactl/internal/pg"
	"github.com/spf13/cobra"
)

var listSort string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pg_hba.conf rules in a table",
	Long:  "Connects to PostgreSQL, discovers pg_hba.conf, parses it, and prints rules in a formatted table. Use --sort to order by column (display only; file order is unchanged).",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listSort, "sort", "", "Sort by column: type, database, user, address, method")
}

func runList(cmd *cobra.Command, _ []string) error {
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

	rules, err := hba.ParseFile(path)
	if err != nil {
		return fmt.Errorf("could not read file (try running with sudo?): %w", err)
	}

	if listSort != "" {
		if !hba.ValidSortColumn(listSort) {
			return fmt.Errorf("invalid --sort %q; use one of: type, database, user, address, method", listSort)
		}
		hba.SortRules(rules, listSort)
	}

	fmt.Printf("File: %s (%d rule(s))\n\n", path, len(rules))
	cli.WriteRulesTable(os.Stdout, rules)
	return nil
}
