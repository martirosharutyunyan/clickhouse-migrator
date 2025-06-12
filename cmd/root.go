/*
Copyright © 2025 Martiros <martiros.harutunyan@gmail.com>

*/
package cmd

import (
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "clickhouse-migrator",
	Short: "clickhouse-migrator for distributed cluster",
	Long: `Clickhouse distributed migrator for distributing goose_db_version_table in cluster.
Redistributing it in case of adding or deleting nodes`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(cmd.PersistentFlags().Lookup("s").Changed)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("dir", cfg.GOOSEMIGRATIONDIR, "migrations directory")
	rootCmd.PersistentFlags().String("table", "goose_db_version", "table name for migrations versioning")
	rootCmd.PersistentFlags().Bool("allow-missing", false, "applies missing (out-of-order) migrations")
	rootCmd.PersistentFlags().String(
		"no-versioning",
		"",
		"apply migration commands with no versioning, in file order, from directory pointed to",
	)
	rootCmd.PersistentFlags().Bool("s", false, "use sequential numbering for new migrations")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable color output (NO_COLOR env variable supported)")
	rootCmd.PersistentFlags().Duration("timeout", 0, "maximum allowed duration for queries to run")
	rootCmd.PersistentFlags().String("v", "", "enable verbose mode")
	rootCmd.PersistentFlags().String("dsn", "", "database connection string")
	rootCmd.PersistentFlags().String("migration-type", "sql", "migration type. supported sql, go")
	rootCmd.PersistentFlags().String("name", "", "migration name")
	rootCmd.PersistentFlags().String("cluster", "", "cluster name")
	rootCmd.PersistentFlags().String("db", "", "database name")
}
