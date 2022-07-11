package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed template/initializer-tmpl.yml
var configTemp string

type configFlag struct {
	appName string
	update  bool
}

var cfgFlag = configFlag{}

func generateConfiguration() {
	if _, err := os.ReadFile(Application); err != nil {
		if t, err := template.New(Application).Parse(configTemp); err != nil {
			fmt.Println(fmt.Sprintf("Failed to load the template, %+v", err))
		} else {
			if f, err := os.Create(Application); err != nil {
				fmt.Println(fmt.Sprintf("Failed to create file, %+v", err))
				return
			} else {
				if err = t.Execute(f, cfgFlag.appName); err != nil {
					fmt.Println(fmt.Sprintf("Failed to create file %v, %+v", Application, err))
				}
				f.Close()
				abs, _ := filepath.Abs(f.Name())
				fmt.Println(fmt.Sprintf("create files: %v successfully", abs))
			}
		}
	} else {
		pwd, _ := os.Getwd()
		fmt.Println(fmt.Sprintf("%v/%v exists", pwd, Application))
		return
	}
}

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "config",
		Short:     "Generate system configuration files application.yml",
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"name", "update"},
		Run: func(cmd *cobra.Command, args []string) {
			generateConfiguration()
		},
	}
	cmd.Flags().StringVarP(&cfgFlag.appName, "name", "n", "gbt-app", "set application name")
	cmd.Flags().BoolVarP(&cfgFlag.update, "update", "u", false, "scan re-generate the configuration")
	return cmd
}
