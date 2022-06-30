/*
Copyright © 2022 ken Cheng <kcheng.mvp@gmail.com>
*/
package cmd

import (
	"context"
	"golang.org/x/mod/modfile"
	"os"
	"strings"

	"github.com/spf13/cobra"
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
	Short:             "Generate project scaffold based on the dependencies",
	Long:              "",
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
	rootCmd.AddCommand(initializerCmd())
}

// examples formats the given examples to the cli.
func examples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}
