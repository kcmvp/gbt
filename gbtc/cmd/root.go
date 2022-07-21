/*
Copyright © 2022 ken Cheng <kcheng.mvp@gmail.com>
*/
package cmd

import (
	"context"
	"github.com/kcmvp/gbt/gbtc/cmd/builder"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
	"github.com/kcmvp/gbt/gbtc/cmd/config"
	"github.com/kcmvp/gbt/gbtc/cmd/hook"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"os"
)

func preValidateE(cmd *cobra.Command, args []string) error {
	if data, err := os.ReadFile("go.mod"); err != nil {
		return common.RunFromRootMsg
	} else {
		if f, err := modfile.Parse("go.mod", data, nil); err != nil {
			return common.RunFromRootMsg
		} else {
			pwd, _ := os.Getwd()
			ctx := context.WithValue(cmd.Context(), common.Mod, f)
			ctx = context.WithValue(ctx, common.ProjectRootDir, pwd)
			cmd.SetContext(ctx)
		}
	}
	return nil
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gbtc",
		PersistentPreRunE: preValidateE,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	cmd.AddCommand(config.NewConfigCmd(), builder.NewBuilderCmd(), hook.NewGitHookCmd())
	return cmd
}
