package cmd

import (
	"fmt"
	"github.com/magefile/mage/mage"
	"github.com/spf13/cobra"
	"os"
)

var mageInitF = func(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(mage.MagefilesDirName); err != nil {
		os.Mkdir(mage.MagefilesDirName, 0755)
		fmt.Printf("create mage folder %s successfully\n", mage.MagefilesDirName)
	} else {
		fmt.Printf("mage folder %s exists\n", mage.MagefilesDirName)
	}
}

func mageInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init mage build files in the root directory",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: mageInitF,
	}
	return cmd
}

func mageBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build current project",
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			mage.Main()
		},
	}
	return cmd
}

func mageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mage [init]",
		Short: "Generate project build file.",
		Args: func(_ *cobra.Command, names []string) error {
			fmt.Println("mage Args")
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("in the mageCMD")
		},
	}
	//cmd.AddCommand(mageInitCmd())
	return cmd
}
