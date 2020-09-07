package main

import (
	"fmt"

	"github.com/nomad-software/findr/cli"
	"github.com/nomad-software/findr/file"
)

func main() {
	options := cli.ParseOptions()

	if options.Help {
		options.PrintUsage()

	} else if options.Valid() {
		file := file.NewWalker(&options)

		err := file.Walk()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
