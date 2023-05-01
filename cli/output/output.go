package output

import (
	"fmt"
)

// Worker is the main output worker.
type Worker struct {
	Queue chan string
	done  chan bool
}

// New creates a new output worker. This worker will print all matches to the
// terminal.
func New() *Worker {
	return &Worker{
		Queue: make(chan string),
		done:  make(chan bool),
	}
}

// StartQueue starts the output queue to process matches.
func (q *Worker) StartQueue() {
	for match := range q.Queue {
		fmt.Println(match)
	}

	q.done <- true
}

// StopQueue stops the output queue.
func (q *Worker) StopQueue() {
	close(q.Queue)
	<-q.done
}
