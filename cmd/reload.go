package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hrodrig/hbactl/internal/pg"
	"github.com/spf13/cobra"
)

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Reload PostgreSQL configuration",
	Long:  "Runs pg_reload_conf() so the server picks up the current pg_hba.conf without restart.",
	RunE:  runReload,
}

func init() {
	rootCmd.AddCommand(reloadCmd)
}

func runReload(cmd *cobra.Command, _ []string) error {
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

	if err := client.ReloadConf(ctx); err != nil {
		return fmt.Errorf("reload failed: %w", err)
	}

	fmt.Fprintln(os.Stdout, "Success: configuration reloaded.")
	return nil
}
