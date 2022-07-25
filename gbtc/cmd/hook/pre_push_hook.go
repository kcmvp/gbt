//go:build gbt

package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kcmvp/gbt/script"
	"log"
	"os"
)

func main() {

	arg1 := os.Args[1]
	arg2 := os.Args[2]
	cqc := script.NewCQC(0.35, 0.85)
	if rep, err := git.PlainOpen(cqc.RootDir()); err == nil {
		if b, err := rep.Branch("hook"); err != nil {
			log.Fatalf("can open branch %v", err)
		} else {
			fmt.Println(b)
		}
		os.Exit(1)
	} else {
		log.Fatalf("failed to open repository %s", cqc.RootDir())
	}
	fmt.Printf("commits are %s, %s", arg1, arg2)
	os.Exit(1)
}
