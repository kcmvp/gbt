package cmd

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/spf13/cobra"
)

func mageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mage",
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

func Build() error {
	return nil
}

func Run() {
	mg.Deps(Build)
	sh.RunV("./start.sh")
}
