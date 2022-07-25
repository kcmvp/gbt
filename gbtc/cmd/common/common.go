package common

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const Mod = "mod"
const Application = "application.yml"
const ApplicationTest = "application-test.yml"

const scriptModule = "github.com/kcmvp/gbt/script"
const ProjectRootDir = "projectRootDir"
const ScriptDir = "scripts"
const VersionedCoverage = "versioned_coverage.data"

var RunFromRootMsg = fmt.Errorf("please run the command from project root")

func GenerateFile(ctx context.Context, content string, name string, data interface{}) {
	pwd := ctx.Value(ProjectRootDir).(string)
	target := filepath.Join(pwd, name)
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create dir %s: %v", dir, err)
	}
	if f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm); err == nil {
		defer f.Close()
		if t, err := template.New(name).Parse(content); err != nil {
			fmt.Println(fmt.Sprintf("Failed to parse template, %+v", err))
		} else {
			if err = t.Execute(f, data); err != nil {
				fmt.Println(fmt.Sprintf("Failed to create file %v, %+v", name, err))
			}
			abs, _ := filepath.Abs(f.Name())
			fmt.Println(fmt.Sprintf("generate files: %v successfully", abs))
		}
	} else {
		if errors.Is(err, os.ErrExist) {
			fmt.Println(fmt.Sprintf("file %s exists", name))
		} else {
			fmt.Println(fmt.Sprintf("failed to generate file %s, %v", name, err))
		}
	}
}

func ImportScript(ctx context.Context, update bool) {
	ImportModule(ctx, scriptModule, update)
}

func ImportModule(ctx context.Context, module string, update bool) {
	mod := ctx.Value(Mod).(*modfile.File)
	has := false
	for _, require := range mod.Require {
		if has = require.Mod.Path == module; has {
			break
		}
	}
	if !has || update {
		action := "installing"
		if has {
			action = "updating"
		}
		fmt.Println(fmt.Sprintf("%s %s", action, module))
		output, _ := exec.Command("go", "get", module).CombinedOutput()
		fmt.Println(string(output))
	}
}
