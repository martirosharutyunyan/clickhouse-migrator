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
}

func OptionsFromCmd(cmd *cobra.Command) []goose.ProviderOption {
	var opts []goose.ProviderOption
	if cmd.Flags().Lookup("v").Changed {
		opts = append(opts, goose.WithVerbose(true))
	}
	if cmd.Flags().Lookup("s").Changed {
		opts = append(opts, goose.WithAllowOutofOrder(true))
	}

	return opts
}
