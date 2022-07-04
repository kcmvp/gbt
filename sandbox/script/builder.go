//go:build gbt

package main

import (
	"fmt"
	"github.com/kcmvp/gbt/script"
)

func main() {
	cqc, _ := script.NewCQC()
	if cqc.Clean() == nil {
		if e := cqc.Test(); e != nil {
			fmt.Println(e)
		}
	}
}
