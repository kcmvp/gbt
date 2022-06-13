package cmd

import (
	"github.com/spf13/cobra"
)

//@todo create application.yml and application-test.yml and import github.com/kcmvp/gbt/env

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init the project configuration file. Run this command at your project root directory",
		Args: func(_ *cobra.Command, names []string) error {
			println("Args: init")
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			println("Run: init")
			//initProject()
		},
		//PreRunE: preValidateE,
	}
	return cmd
}
