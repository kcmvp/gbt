/*
Copyright © 2022 ken Cheng <kcheng.mvp@gmail.com>
*/
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"os"
)

func preValidateE(cmd *cobra.Command, args []string) error {
	if data, err := os.ReadFile("go.mod"); err != nil {
		return NOT_IN_ROOT
	} else {
		if f, err := modfile.Parse("go.mod", data, nil); err != nil {
			return NOT_IN_ROOT
		} else {
			context.WithValue(cmd.Context(), Mod, f)
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:               "gbt-cli",
	Short:             "gbt scaffold commands",
	Long:              "Please run this command in the project root directory",
	PersistentPreRunE: preValidateE,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(builderCmd())
}
