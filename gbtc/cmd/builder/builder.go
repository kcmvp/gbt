package builder

import (
	"context"
	_ "embed"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
	"github.com/spf13/cobra"
)

//go:embed builder.tmpl
var builderTmp string

type builderFlag struct {
	update bool
}

func (f *builderFlag) Update() bool {
	return f.update
}

var bFlag = builderFlag{}

func generateBuilder(ctx context.Context) {
	common.GenerateFile(ctx, builderTmp, "scripts/builder.go", nil)
	common.ImportScript(ctx, bFlag.update)
}

func NewBuilderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "builder",
		Short:     "Generate build scripts for the project",
		Long:      "Generate script/builder.go at project root, you can build project by execute : go run script/builder.go",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			generateBuilder(cmd.Context())
		},
	}
	cmd.Flags().BoolVarP(&bFlag.update, "update", "u", false, "update to the latest github.com/kcmvp/gbt/script@latest")
	return cmd
}
