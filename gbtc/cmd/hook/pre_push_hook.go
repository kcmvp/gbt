//go:build gbt

package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/kcmvp/gbt/script"
	"os"
)

func main() {
	git.PlainOpen("")
	cqc := script.NewCQC()
	cqc.Clean().Test()
	if cqc.Error() != nil {
		os.Exit(1)
	} else {
		// scripts/coverage.data must be exists
		// todo compare coverage with the test output must be equal
		// todo local coverage must >= remote coverage
		os.Exit(0)
	}
}
