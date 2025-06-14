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

// downToCmd represents the downTo command
var downToCmd = &cobra.Command{
	Use:   "down-to",
	Short: "down to specified version",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("must specify version")
		}
		cfg.InitGoose(cmd)

		conf, err := cfg.NewConfig(cmd)
		if err != nil {
			return err
		}

		provider, err := database.NewProvider(conf, cfg.OptionsFromCmd(cmd)...)
		if err != nil {
			return err
		}

		version, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		res, err := provider.DownTo(cmd.Context(), int64(version))
		if err != nil {
			return err
		}

		fmt.Println("Roll backed following migrations")
		for _, row := range res {
			fmt.Println(row.Source.Path)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(downToCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downToCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downToCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
