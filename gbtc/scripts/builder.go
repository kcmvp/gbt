//go:build gbt




package main

import "github.com/kcmvp/gbt/script"

func main() {
	cqc := script.NewCQC()
	cqc.Clean().Test().Build()
}
