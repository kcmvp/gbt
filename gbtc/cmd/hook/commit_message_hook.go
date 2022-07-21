//go:build gbt

package main

import (
	"fmt"
	"os"
	"regexp"
)

const MsgPattern = "#[0-9]{1,7}:.*"
const Separator = "#[0-9]{1,7}:"
const MinLength = 10

func main() {
	input, _ := os.ReadFile(os.Args[1])
	msg := checkFormat(string(input))
	checkLength(msg)
	os.Exit(0)
}

func checkFormat(msg string) string {
	reg, err := regexp.Compile(MsgPattern)
	sp, _ := regexp.Compile(Separator)
	if err != nil {
		fmt.Println(fmt.Sprintf("internal error %v", err))
		os.Exit(1)
	}
	if !reg.MatchString(msg) {
		fmt.Println("commit message must follow format #{number}: xxxxxx")
		os.Exit(1)
	}
	items := sp.Split(msg, -1)
	return items[1]
}

func checkLength(msg string) {
	if len(msg) < MinLength {
		fmt.Println(fmt.Sprintf("commit message is at least %d characters", MinLength))
		os.Exit(1)
	}
}
