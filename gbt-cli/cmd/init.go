package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

//@todo create application.yml and application-test.yml and import github.com/kcmvp/gbt/env

const application = "application.yml"
const applicationTest = "application-test.yml"

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init the project configuration file. Run this command at your project root directory",
		Args: func(_ *cobra.Command, names []string) error {
			println("Args: init")
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			initBuilder(cmd, names)
		},
	}
	return cmd
}

func initBuilder(cmd *cobra.Command, names []string) {
	for _, f := range []string{application, applicationTest} {
		if _, err := os.Stat(f); err != nil {
			// init the file
		}
	}
}
