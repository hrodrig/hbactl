package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var connStr string
var hbaFilePath string

// Version is set at build time via ldflags (e.g. -X github.com/hrodrig/hbactl/cmd.Version=v0.1.0). Default: "dev".
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "hbactl",
	Short: "Manage PostgreSQL pg_hba.conf safely",
	Long:  "hbactl is a CLI to manage and validate PostgreSQL Host-Based Authentication (pg_hba.conf) with safety and ease.",
}

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("hbactl {{.Version}}\n")
	rootCmd.PersistentFlags().StringVarP(&connStr, "conn", "c", "", "PostgreSQL connection string (default: DATABASE_URL env)")
	rootCmd.PersistentFlags().StringVarP(&hbaFilePath, "file", "f", "", "Path to pg_hba.conf (if set, list uses it and may skip connection; for multiple servers, pass path per run)")
}

// connString returns the connection string: flag if set, else DATABASE_URL.
func connString() string {
	if connStr != "" {
		return connStr
	}
	return os.Getenv("DATABASE_URL")
}

// filePath returns the explicit pg_hba.conf path from --file / -f (empty if not set).
func filePath() string {
	return hbaFilePath
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}
