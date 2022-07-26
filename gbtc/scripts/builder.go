//go:build gbt

package main

import (
	"flag"
	"github.com/kcmvp/gbt/gbtc/cmd/common"
	"github.com/kcmvp/gbt/script"
)

func main() {
	build := true
	flag.BoolVar(&build, "build", true, "build project or not")
	cqc := script.NewCQC().WithFlag(common.FlagBuild, build)
	cqc.Clean().Test().Build()
}
