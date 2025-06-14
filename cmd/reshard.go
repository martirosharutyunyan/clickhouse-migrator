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

// reshardCmd represents the reshard command
var reshardCmd = &cobra.Command{
	Use:   "reshard",
	Short: "reshard the migrations table",
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

		res, err := database.Reshard(cmd.Context(), conf, provider)
		if err != nil {
			return err
		}

		fmt.Println("Resharded following migrations")
		for _, row := range res {
			fmt.Printf("source: %s, shard_num %d\n", row.Source, row.ShardNumber)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(reshardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reshardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reshardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
