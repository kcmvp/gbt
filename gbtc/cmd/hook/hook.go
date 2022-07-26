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
	if root, err := common.ProjectRoot(pwd); err == nil {
		scriptDir := filepath.Join(root, common.ScriptDir)
		os.MkdirAll(scriptDir, os.ModePerm)
		for k, v := range hookMap {
			hook := filepath.Join(root, ".git", "hooks", k)
			if f, err := os.OpenFile(hook, os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm); err == nil {
				fmt.Println(fmt.Sprintf("generate %s hook", k))
				f.WriteString("#!/bin/sh\n\n")
				f.WriteString(fmt.Sprintf("go run %s $1 $2\n", filepath.Join(scriptDir, v)))
				f.Close()
			} else if errors.Is(err, os.ErrExist) {
				fmt.Println(fmt.Sprintf("%s exists", hook))
			}
			// generate the hook when it does not exist
			tn := strings.Replace(v, ".go", ".tmpl", 1)
			if data, err := templateDir.ReadFile(tn); err == nil {
				common.GenerateFile(ctx, string(data), filepath.Join(scriptDir, v), nil)
			}
		}
		common.ImportScript(ctx, false)
		common.ImportModule(ctx, goGit, false)
	} else {
		fmt.Println(err)
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
