package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/nomad-software/findr/cli"
)

// Walker is the main file walker, it coordinates matching options and the path queue.
type Walker struct {
	Group   sync.WaitGroup
	Ignore  *regexp.Regexp
	Options *cli.Options
	Output  *cli.Output
	Pattern *regexp.Regexp
}

// NewWalker creates a new file walker.
func NewWalker(opt *cli.Options) Walker {
	var w Walker

	w.Options = opt
	w.Pattern = w.compile(opt.Regex)

	if w.Options.Ignore != "" {
		w.Ignore = w.compile(w.Options.Ignore)
	}

	w.Output = &cli.Output{
		Console: make(chan string),
		Closed:  make(chan bool),
	}

	return w
}

// Walk starts walking through the directory specified in the options and starts
// processing any matched files.
func (w *Walker) Walk() error {
	go w.Output.Start()

	err := filepath.Walk(w.Options.Dir, func(fullPath string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		w.Group.Add(1)
		go w.matchPath(fullPath)

		return nil
	})

	w.Group.Wait()
	w.Output.Stop()

	return err
}

// matchPath matches paths and passed them to be output.
func (w *Walker) matchPath(fullPath string) {
	defer w.Group.Done()

	if w.Ignore != nil && w.Ignore.MatchString(fullPath) {
		return
	}

	matched, err := filepath.Match(w.Options.Glob, path.Base(fullPath))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if matched && w.Pattern.MatchString(fullPath) {
		w.Output.Console <- fullPath
	}
}

// Check that a regex pattern compiles.
func (w *Walker) compile(pattern string) (regex *regexp.Regexp) {
	if w.Options.Case {
		regex, _ = regexp.Compile(pattern)
	} else {
		regex, _ = regexp.Compile("(?i)" + pattern)
	}

	return regex
}
