/*
Copyright © 2022 ken Cheng <kcheng.mvp@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"os"
)

func preValidateE(cmd *cobra.Command, args []string) error {
	if data, err := os.ReadFile("go.mod"); err != nil {
		return runFromRootMsg
	} else {
		if _, err = modfile.Parse("go.mod", data, nil); err != nil {
			return runFromRootMsg
		}
	}
	return nil
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gbtc",
		PersistentPreRunE: preValidateE,
		//Run: func(cmd *cobra.Command, args []string) {
		//	cmd.Usage()
		//},
	}
	cmd.AddCommand(generator())
	return cmd
}

func generator() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gen",
		Short:             "generate project scaffold",
		PersistentPreRunE: preValidateE,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	cmd.AddCommand(ConfigCmd(), BuilderCmd())
	return cmd
}
