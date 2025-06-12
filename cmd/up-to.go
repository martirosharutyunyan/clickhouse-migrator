/*
Copyright © 2025 Martiros <martiros.harutunyan@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/cfg"
	"github.com/martirosharutyunyan/clickhouse-migrator/pkg/database"
	"strconv"

	"github.com/spf13/cobra"
)

// upToCmd represents the upTo command
var upToCmd = &cobra.Command{
	Use:   "up-to",
	Short: "up migrations to specified version",
	Long:  ``,
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

		version, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		res, err := provider.UpTo(cmd.Context(), int64(version))
		if err != nil {
			return err
		}

		fmt.Println(res)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upToCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upToCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upToCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
