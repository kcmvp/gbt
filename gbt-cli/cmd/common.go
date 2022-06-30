package cmd

import "fmt"

const Mod = "mod"
const Application = "application.yml"
const ApplicationTest = "application-test.yml"

var NOT_IN_ROOT = fmt.Errorf("please run the command in the root directory")
