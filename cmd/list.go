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
var listNoIndex bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List pg_hba.conf rules in a table",
	Long:  "Connects to PostgreSQL, discovers pg_hba.conf, parses it, and prints rules in a formatted table. Use --sort to order by column (display only; file order is unchanged). Use --group-by user to print separators between users. Use --no-index to omit the rule index column for copy-paste friendly output.",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&listSort, "sort", "", "Sort by column: type, database, user, address, method")
	listCmd.Flags().StringVar(&listGroupBy, "group-by", "", "Print visual separators by column (e.g. user); implies --sort by that column if not set")
	listCmd.Flags().BoolVar(&listNoIndex, "no-index", false, "Omit the # column so output is copy-paste friendly (no rule numbers)")
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

	rwl, err := hba.ParseFileWithLineNumbers(path)
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
		hba.SortRulesWithLine(rwl, sortCol)
	}

	fmt.Printf("File: %s (%d rule(s))\n\n", path, len(rwl))
	if listNoIndex {
		if listGroupBy == "user" {
			cli.WriteRulesTableGroupedByUserNoIndex(os.Stdout, rwl)
		} else {
			cli.WriteRulesTableNoIndex(os.Stdout, rwl)
		}
	} else {
		if listGroupBy == "user" {
			cli.WriteRulesTableGroupedByUser(os.Stdout, rwl)
		} else {
			cli.WriteRulesTableWithIndex(os.Stdout, rwl)
		}
	}
	return nil
}
