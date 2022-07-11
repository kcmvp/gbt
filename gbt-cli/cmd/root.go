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
		//else {
		//	context.WithValue(cmd.Context(), Mod, f)
		//}
	}
	return nil
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gbt-cli",
		Short:             "gbt scaffold commands",
		PersistentPreRunE: preValidateE,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	cmd.AddCommand(ConfigCmd())
	cmd.AddCommand(BuilderCmd())
	return cmd
}
