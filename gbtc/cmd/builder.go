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
	"text/template"
)

//go:embed template/builder.tmpl
var builderTmp string

var defaultBuilderDir = "scripts"
var defaultBuilder = "builder.go"
var scriptModule = "github.com/kcmvp/gbt/script"

type builderFlag struct {
	update bool
}

var bFlag = builderFlag{}

func importScriptModule(ctxt context.Context) {
	data, _ := os.ReadFile("go.mod")
	mod, _ := modfile.Parse("go.mod", data, nil)
	has := false
	for _, require := range mod.Require {
		if has = require.Mod.Path == scriptModule; has {
			break
		}
	}
	if !has || bFlag.update {
		action := "installing"
		if has {
			action = "updating"
		}
		fmt.Println(fmt.Sprintf("%s %s", action, scriptModule))
		output, _ := exec.Command("go", "get", scriptModule).CombinedOutput()
		fmt.Println(string(output))
	}
}

func createBuilder() {
	builder := fmt.Sprintf("%s/%s", defaultBuilderDir, defaultBuilder)
	if _, err := os.Stat(builder); err != nil {
		os.Mkdir(defaultBuilderDir, os.ModePerm)
		if t, err := template.New("defaultBuilder").Parse(builderTmp); err != nil {
			fmt.Println(fmt.Sprintf("Failed to load the template, %+v", err))
		} else {
			if f, err := os.Create(builder); err != nil {
				fmt.Println(fmt.Sprintf("Failed to create file, %+v", err))
				return
			} else {
				if err = t.Execute(f, nil); err != nil {
					fmt.Println(fmt.Sprintf("Failed to create file %v, %+v", Application, err))
				}
				f.Close()
				abs, _ := filepath.Abs(f.Name())
				fmt.Println(fmt.Sprintf("create files: %v successfully", abs))
			}
		}
	}
}

func generateBuilder(ctxt context.Context) {
	if _, err := os.Stat(defaultBuilderDir); err != nil {
		fmt.Println("Creating directory: scripts")
		if err = os.Mkdir(defaultBuilderDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create defaultBuilderDir %s: %v", defaultBuilderDir, err)
		}
	}
	f := filepath.Join(defaultBuilderDir, defaultBuilder)
	if _, err := os.Stat(f); err == nil {
		fmt.Println(fmt.Sprintf("File %s exists", f))
		return
	}
	script.InstallDependencies()
	createBuilder()
	importScriptModule(ctxt)
}

func BuilderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "builder",
		Short:     "Generate build script for the project",
		Long:      "Generate script/defaultBuilder.go at project root, you can build project by execute : go run script/defaultBuilder.go",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			generateBuilder(cmd.Context())
		},
	}
	cmd.Flags().BoolVarP(&bFlag.update, "update", "u", false, "update to the latest github.com/kcmvp/gbt/script@latest")
	return cmd
}
