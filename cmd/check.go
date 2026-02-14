package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hrodrig/hbactl/internal/pg"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate pg_hba.conf for syntax errors",
	Long:  "Queries pg_hba_file_rules to detect any lines PostgreSQL could not parse. Exit 0 if OK, 1 if errors found.",
	RunE:  runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, _ []string) error {
	conn := connString()
	if conn == "" {
		return fmt.Errorf("no connection: set DATABASE_URL or use --conn")
	}

	ctx := context.Background()
	client, err := pg.NewClient(ctx, conn)
	if err != nil {
		return fmt.Errorf("could not connect to PostgreSQL: %w", err)
	}
	defer client.Close()

	errs, err := client.HBAFileErrors(ctx)
	if err != nil {
		return fmt.Errorf("could not read pg_hba_file_rules: %w", err)
	}

	if len(errs) == 0 {
		fmt.Fprintln(os.Stdout, "OK: no syntax errors in pg_hba.conf")
		return nil
	}

	fmt.Fprintln(os.Stderr, "Error: syntax errors in pg_hba.conf:")
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "  line %d: %s\n", e.LineNumber, e.Error)
	}
	// Exit 1 will be set by main when we return a non-nil error. So we need to return an error.
	return fmt.Errorf("%d syntax error(s) found", len(errs))
}
