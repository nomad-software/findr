package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is a color friendly pipe.
	Stdout = colorable.NewColorableStdout()

	// Stderr is a color friendly pipe.
	Stderr = colorable.NewColorableStderr()
)

type Output struct {
	Closed  chan bool
	Console chan string
}

func (this *Output) Start() {
	var total int
	for path := range this.Console {
		total++
		color.Magenta(path)
	}
	fmt.Printf("Found %d file(s)\n", total)
	this.Closed <- true
}

func (this *Output) Stop() {
	close(this.Console)
	<-this.Closed
}
