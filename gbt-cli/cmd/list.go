package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "init the project configuration file. Run this command at your project root directory",
		Args: func(_ *cobra.Command, names []string) error {
			fmt.Println("list:Args")
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			fmt.Println("list:Run")
		},
	}
	return cmd
}
