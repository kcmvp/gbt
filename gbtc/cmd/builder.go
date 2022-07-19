package cmd

import (
	"context"
	_ "embed"
	"fmt"
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

var builderDir = "scripts"
var builderScript = "builder.go"
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
	builder := fmt.Sprintf("%s/%s", builderDir, builderScript)
	if _, err := os.Stat(builder); err != nil {
		os.Mkdir(builderDir, os.ModePerm)
		if t, err := template.New("builderScript").Parse(builderTmp); err != nil {
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
	if _, err := os.Stat(builderDir); err != nil {
		fmt.Println("Creating directory: scripts")
		if err = os.Mkdir(builderDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create builderDir %s: %v", builderDir, err)
		}
	}
	f := filepath.Join(builderDir, builderScript)
	if _, err := os.Stat(f); err == nil {
		fmt.Println(fmt.Sprintf("File %s exists", f))
		return
	}
	createBuilder()
	importScriptModule(ctxt)
}

func BuilderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "builder",
		Short:     "Generate build scripts for the project",
		Long:      "Generate script/builderScript.go at project root, you can build project by execute : go run script/builderScript.go",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"update"},
		Run: func(cmd *cobra.Command, args []string) {
			generateBuilder(cmd.Context())
		},
	}
	cmd.Flags().BoolVarP(&bFlag.update, "update", "u", false, "update to the latest github.com/kcmvp/gbt/script@latest")
	return cmd
}
