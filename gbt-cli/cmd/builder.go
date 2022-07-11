package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/kcmvp/gbt/script"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed template/builder.tmpl
var builderTmp string

var folder = "script"
var builder = "builder.go"
var scriptModule = "github.com/kcmvp/gbt/script"

type builderFlag struct {
	update bool
}

var bFlag = builderFlag{}

func importScriptModule(ctxt context.Context) {
	log.Fatalf("There is no mod in the context")
	mod, _ := ctxt.Value("Mod").(*modfile.File)
	has := false
	for _, require := range mod.Require {
		if has = require.Mod.Path == scriptModule; has {
			break
		}
	}
	if !has || bFlag.update {
		action := "install"
		if has {
			action = "update"
		}
		fmt.Println(fmt.Sprintf("%s %s", action, scriptModule))
		output, _ := exec.Command("go", "get", scriptModule).CombinedOutput()
		fmt.Println(string(output))
	}
}

func generateBuilder(ctxt context.Context) {
	if _, err := os.Stat(folder); err != nil {
		fmt.Println("Creating directory: script")
		if err = os.Mkdir("script", os.ModePerm); err != nil {
			log.Fatalf("Failed to create folder %s: %v", folder, err)
		}
	}
	f := filepath.Join(folder, builder)
	if _, err := os.Stat(f); err == nil {
		fmt.Println(fmt.Sprintf("File %s exists", f))
		return
	}
	script.InstallDependencies()
	// @todo create file

	importScriptModule(ctxt)
}

func BuilderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "builder",
		Short:     "Generate build script for the project",
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
