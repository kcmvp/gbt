package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func generateCmd() *cobra.Command {
	//var target string
	cmd := &cobra.Command{
		Use:   "generate xxxx",
		Short: "initialize an environment with zero or more schemas",
		Example: examples(
			"ent init Example",
			"ent init --target entv1/schema User Group",
		),
		Args: func(_ *cobra.Command, names []string) error {
			fmt.Println("Args: generate")
			return nil
		},
		Run: func(cmd *cobra.Command, names []string) {
			fmt.Println("Run: generate")
		},
	}
	//cmd.Flags().StringVar(&target, "target", defaultSchema, "target directory for schemas")
	return cmd
}
