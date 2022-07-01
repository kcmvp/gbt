//go:build gbt

package main

import (
	"fmt"
	"github.com/kcmvp/gbt/script"
)

func main() {
	b, _ := script.DefaultBuilder()
	if b.Clean() == nil {
		if e := b.Test(); e != nil {
			fmt.Println(e)
		}
	}
}
