//go:build gbt

package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])
	//os.Exit(1)

	//cqc := script.NewCQC(0.35, 0.85)
	//cqc.Clean().Test()
	//if cqc.Error() != nil {
	//	log.Fatalf("%v", cqc.Error())
	//} else {
	//	if rep, err := git.PlainOpen(cqc.RootDir()); err == nil {
	//		rep.Head()
	//		rep.Prune()
	//		commit, _ := object.GetCommit()
	//		commit.Patch()
	//	} else {
	//		log.Fatalf("can't detect the git repository %s", cqc.RootDir())
	//	}
	//	os.Exit(0)
	//}
}
