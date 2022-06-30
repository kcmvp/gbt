package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed initializer-tmpl.yml
var tmp string

var appName string

func defaultApplication() {
	if _, err := os.ReadFile(Application); err != nil {
		if t, err := template.New(Application).Parse(tmp); err != nil {
			fmt.Println(fmt.Sprintf("Failed to load the template, %+v", err))
		} else {
			if f, err := os.Create(Application); err != nil {
				fmt.Println(fmt.Sprintf("Failed to create file, %+v", err))
				return
			} else {
				if err = t.Execute(f, appName); err != nil {
					fmt.Println(fmt.Sprintf("Failed to create file %v, %+v", Application, err))
				}
				f.Close()
				abs, _ := filepath.Abs(filepath.Dir(f.Name()))
				fmt.Println(fmt.Sprintf("create files: %+v successfully", abs))
			}
		}
	} else {
		fmt.Println(fmt.Sprintf("%v exists", Application))
		return
	}
}

func initializerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "init",
		Short:     "Initialize system configuration files: application.yml",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"name"},
		Run: func(cmd *cobra.Command, args []string) {
			defaultApplication()
		},
	}
	cmd.Flags().StringVarP(&appName, "name", "n", "gbt-app", "Set application name")
	return cmd
}
