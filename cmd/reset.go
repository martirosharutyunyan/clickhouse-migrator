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

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all migrations",
	Long:  ``,
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

		res, err := database.Reset(cmd.Context(), provider)
		if err != nil {
			return err
		}

		if res == nil {
			return nil
		}

		fmt.Println("Reseted migrations")
		for _, row := range res {
			fmt.Println(row.Source.Path)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
