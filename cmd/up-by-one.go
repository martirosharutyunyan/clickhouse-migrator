/*
Copyright © 2025 Martiros <martiros.harutunyan@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/database"
	"github.com/spf13/cobra"
)

// upByOneCmd represents the upByOne command
var upByOneCmd = &cobra.Command{
	Use:   "up-by-one",
	Short: "up migrations one by one",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg.InitGoose(cmd)

		conf, err := cfg.NewConfig(cmd)
		if err != nil {
			return err
		}

		provider, err := database.NewProvider(conf, cfg.OptionsFromCmd(cmd)...)
		if err != nil {
			return err
		}

		res, err := provider.UpByOne(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Printf("Applied migration: %s", res.Source.Path)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upByOneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upByOneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upByOneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
