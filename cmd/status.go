/*
Copyright © 2025 Martiros <martiros.harutunyan@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/database"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.InitGoose(cmd)

		conf, err := cfg.NewConfig(cmd)
		if err != nil {
			return err
		}

		provider, err := database.NewProvider(conf)
		if err != nil {
			return err
		}

		res, err := provider.Status(cmd.Context())
		if err != nil {
			return err
		}
		fmt.Println("status: ", lo.Map(res, func(item *goose.MigrationStatus, index int) Version {
			return Version{
				Version:   fmt.Sprintf("%d_%s", item.Source.Version, item.Source.Path),
				IsApplied: string(item.State),
			}
		}))

		return nil
	},
}

type Version struct {
	Version   string
	IsApplied string
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
