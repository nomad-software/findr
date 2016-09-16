package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fatih/color"
	"github.com/nomad-software/findr/cli"
)

type Handler struct {
	Group   sync.WaitGroup
	Ignore  *regexp.Regexp
	Options *cli.Options
	Output  *cli.Output
	Pattern *regexp.Regexp
}

func NewHandler(options *cli.Options) Handler {
	var handler Handler

	handler.Options = options
	handler.Pattern = handler.compile(options.Pattern)

	if handler.Options.Ignore != "" {
		handler.Ignore = handler.compile(handler.Options.Ignore)
	}

	handler.Output = &cli.Output{
		Console: make(chan string),
		Closed:  make(chan bool),
	}

	return handler
}

func (this *Handler) Walk() error {
	go this.Output.Start()

	err := filepath.Walk(this.Options.Dir, func(fullPath string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		this.Group.Add(1)
		go this.matchPath(fullPath)

		return nil
	})

	this.Group.Wait()
	this.Output.Stop()

	return err
}

func (this *Handler) matchPath(fullPath string) {
	defer this.Group.Done()

	if this.Ignore != nil && this.Ignore.MatchString(fullPath) {
		return
	}

	matched, err := filepath.Match(this.Options.File, path.Base(fullPath))
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(err.Error()))
		return
	}

	if matched && this.Pattern.MatchString(fullPath) {
		this.Output.Console <- fullPath
	}
}

func (this *Handler) compile(pattern string) (regex *regexp.Regexp) {
	if this.Options.Case {
		regex, _ = regexp.Compile(pattern)
	} else {
		regex, _ = regexp.Compile("(?i)" + pattern)
	}

	return regex
}
