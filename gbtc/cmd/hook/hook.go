package hook

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

//go:embed *.tmpl
var templateDir embed.FS

var hookMap = map[string]string{"commit-msg": "commit_message_hook.go",
	"pre-push": "pre_push_hook.go"}

var goGit = "github.com/go-git/go-git/v5"

func processHookScript(ctx context.Context) {
	pwd, _ := os.Getwd()
	lookup := []string{pwd, filepath.Dir(pwd)}
	gitDir := ""
	for _, path := range lookup {
		gitDir = filepath.Join(path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			break
		} else {
			gitDir = ""
		}
	}
	if len(gitDir) > 0 {
		for k, v := range hookMap {
			script := filepath.Join(gitDir, "hooks", k)
			f, err := os.OpenFile(script, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
			if err == nil {
				fmt.Println(fmt.Sprintf("generate %s hook", k))
				f.WriteString("#!/bin/sh\n\n")
				f.WriteString(fmt.Sprintf("go run %s $1\n", filepath.Join(pwd, common.ScriptDir, v)))
				f.Close()
			} else if errors.Is(err, os.ErrExist) {
				fmt.Println(fmt.Sprintf("%s exists", script))
			}
			// generate the hook when it does not exist
			tn := strings.Replace(v, ".go", ".tmpl", 1)
			if data, err := templateDir.ReadFile(tn); err == nil {
				common.GenerateFile(ctx, string(data), filepath.Join(common.ScriptDir, v), nil)
			}
		}
		common.ImportScript(ctx, false)
		common.ImportModule(ctx, goGit, false)
	} else {
		fmt.Printf("%s is no git repository", pwd)
	}
}

func NewGitHookCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "githook",
		Short: "Generate git local hooks for project",
		Run: func(cmd *cobra.Command, args []string) {
			processHookScript(cmd.Context())
		},
	}
	return cmd
}
