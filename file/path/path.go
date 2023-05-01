package path

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/nomad-software/findr/cli/output"
	"github.com/nomad-software/findr/sync"
)

// Worker is the main path worker.
type Worker struct {
	Queue  chan Match
	Output *output.Worker
	done   chan bool
}

// Match wraps a file path and the pattern being matched agasint it.
type Match struct {
	FullPath string
	Ignore   *regexp.Regexp
	Regex    *regexp.Regexp
	Glob     string
}

// New creates a new path worker. This worker will match paths against the
// specified options and if matched, pass them to the output worker for printing
// to the terminal.
func New() *Worker {
	return &Worker{
		Queue:  make(chan Match),
		Output: output.New(),
		done:   make(chan bool),
	}
}

// StartQueue starts the path queue to process matches.
func (q *Worker) StartQueue() {
	go q.Output.StartQueue()

	sync.CreateWorkers(q.matchPaths, 100)

	q.done <- true
}

// StopQueue stops the path queue.
func (q *Worker) StopQueue() {
	close(q.Queue)
	<-q.done

	q.Output.StopQueue()
}

// matchPaths matches file paths for outputing.
func (q *Worker) matchPaths() error {
	for work := range q.Queue {
		if work.Ignore != nil && work.Ignore.MatchString(work.FullPath) {
			continue
		}

		matched, err := filepath.Match(work.Glob, path.Base(work.FullPath))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		if matched {
			q.Output.Queue <- work.FullPath
		}
	}

	return nil
}
