package cmd

import (
	"fmt"
	"github.com/rogpeppe/go-internal/modfile"
	"log"
	"os"
)

func parseMod() error {
	data, err := os.ReadFile("go.mod");
	if err != nil {
		log.Fatalln("Can't find go.mod file in current folder")
	}
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Go.Version)
	return err
}
