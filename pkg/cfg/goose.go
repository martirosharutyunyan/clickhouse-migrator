package cfg

import (
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

func InitGoose(cmd *cobra.Command) {
	if cmd.Flags().Lookup("v").Changed {
		goose.SetVerbose(true)
	}
	if cmd.Flags().Lookup("s").Changed {
		goose.SetSequential(true)
	}
	goose.SetTableName(cmd.Flags().Lookup("table").Value.String())
}
