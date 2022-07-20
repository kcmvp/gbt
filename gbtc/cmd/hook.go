package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var hookMap = map[string]string{"commit-msg": "commit_message_hook.go",
	"pre-push": "pre_push_hook.go"}

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
			f, err := os.OpenFile(filepath.Join(gitDir, "hooks", k), os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
			if err == nil {
				hook := filepath.Join(pwd, fmt.Sprintf("scripts/%s", v))
				f.WriteString("#!/bin/sh\n\n")
				f.WriteString(fmt.Sprintf("go run %s $1\n", hook))
				f.Close()
			}
		}
	} else {
		fmt.Printf("%s is no git repository", pwd)
	}
}

func GitHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "githook",
		Short: "Generate git local hooks for project",
		Run: func(cmd *cobra.Command, args []string) {
			processHookScript(cmd.Context())
		},
	}
	return cmd
}
