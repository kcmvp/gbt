package cmd

import (
	"context"
	"fmt"
	"github.com/rogpeppe/go-internal/modfile"
	"github.com/spf13/cobra"
	"os"
)

const mod = "mod"

var not_in_root = fmt.Errorf("please run the command in the root directory")

func preValidateE(cmd *cobra.Command, args []string) error {
	if data, err := os.ReadFile("go.mod"); err != nil {
		return not_in_root
	} else {
		if f, err := modfile.Parse("go.mod", data, nil); err != nil {
			return not_in_root
		} else {
			context.WithValue(cmd.Context(), mod, f)
		}
	}
	return nil
}
