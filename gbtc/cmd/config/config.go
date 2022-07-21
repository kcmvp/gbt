package config

import (
	_ "embed"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
	"github.com/spf13/cobra"
)

//go:embed initializer-tmpl.yml
var configTemp string

type configFlag struct {
	appName string
	update  bool
}

var cfgFlag = configFlag{}

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "config",
		Short:     "Generate system configuration files application.yml",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"name", "update"},
		Run: func(cmd *cobra.Command, args []string) {
			common.GenerateFile(cmd.Context(), configTemp, common.Application, cfgFlag.appName)
		},
	}
	cmd.Flags().StringVarP(&cfgFlag.appName, "name", "n", "gbt-app", "set application name")
	cmd.Flags().BoolVarP(&cfgFlag.update, "update", "u", false, "scan re-generate the configuration")
	return cmd
}
