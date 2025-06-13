/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/database"

	"github.com/spf13/cobra"
)

// rerunCmd represents the reshard command
var rerunCmd = &cobra.Command{
	Use:   "rerun",
	Short: "it resets and runs migrations",
	Long:  `Can be used to reshard migrations table in case of adding or deleting nodes`,
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

		res, err := database.Reshard(cmd.Context(), provider)
		if err != nil {
			return err
		}

		fmt.Println("Rerun following migrations")
		for _, row := range res {
			fmt.Println(row.Source.Path)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rerunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reshardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reshardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
