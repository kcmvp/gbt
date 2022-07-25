//go:build gbt

package main

import (
	"github.com/go-git/go-git/v5"
	"github.com/kcmvp/gbt/script"
	"log"
	"os"
)

func main() {
	git.PlainOpen("")
	cqc := script.NewCQC(0.35, 0.85)
	cqc.Clean().Test()
	if cqc.Error() != nil {
		log.Fatalf("%v", cqc.Error())
	} else {
		// scripts/coverage.data must be exists
		// todo compare coverage with the test output must be equal
		// todo local coverage must >= remote coverage
		// git diff HEAD^^ gbtc/scripts/coverage.data
		os.Exit(0)
	}
}
