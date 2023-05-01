package main

import (
	"flag"
	"fmt"

	"github.com/nomad-software/findr/cli"
	"github.com/nomad-software/findr/file/walker"
)

func main() {
	options := cli.ParseOptions()

	if options.Help {
		flag.Usage()

	} else if options.Valid() {
		err := walker.New(options).Walk()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
