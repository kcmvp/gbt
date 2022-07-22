//go:build gbt

package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/kcmvp/gbt/script"
)

func main() {
	git.PlainOpen("")
	cqc := script.NewCQC()
	cqc.Clean().Test()
}
