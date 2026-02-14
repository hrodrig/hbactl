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
var listGroupBy string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pg_hba.conf rules in a table",
	Long:  "Connects to PostgreSQL, discovers pg_hba.conf, parses it, and prints rules in a formatted table. Use --sort to order by column (display only; file order is unchanged). Use --group-by user to print separators between users.",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listSort, "sort", "", "Sort by column: type, database, user, address, method")
	listCmd.Flags().StringVar(&listGroupBy, "group-by", "", "Print visual separators by column (e.g. user); implies --sort by that column if not set")
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

	sortCol := listSort
	if listGroupBy != "" {
		if !hba.ValidSortColumn(listGroupBy) {
			return fmt.Errorf("invalid --group-by %q; use one of: type, database, user, address, method", listGroupBy)
		}
		if sortCol == "" {
			sortCol = listGroupBy
		}
	}
	if sortCol != "" {
		if !hba.ValidSortColumn(sortCol) {
			return fmt.Errorf("invalid --sort %q; use one of: type, database, user, address, method", sortCol)
		}
		hba.SortRules(rules, sortCol)
	}

	fmt.Printf("File: %s (%d rule(s))\n\n", path, len(rules))
	if listGroupBy == "user" {
		cli.WriteRulesTableGroupedByUser(os.Stdout, rules)
	} else {
		cli.WriteRulesTable(os.Stdout, rules)
	}
	return nil
}
