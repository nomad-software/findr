package cli

import (
	"fmt"
)

// Output holds the output channels of the work queue.
type Output struct {
	Closed  chan bool
	Console chan string
}

// Start the output channel.
func (opt *Output) Start() {
	for path := range opt.Console {
		fmt.Println(path)
	}
	opt.Closed <- true
}

// Stop the output channel.
func (opt *Output) Stop() {
	close(opt.Console)
	<-opt.Closed
}
