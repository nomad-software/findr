package file

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/nomad-software/findr/cli"
)

type Handler struct {
	Options *cli.Options
	Group   sync.WaitGroup
	Output  *cli.Output
	Pattern *regexp.Regexp
	Ignore  *regexp.Regexp
}

func (this *Handler) Init(options *cli.Options) {
	this.Options = options
	this.Pattern = this.compile(options.Pattern)

	if this.Options.Ignore != "" {
		this.Ignore = this.compile(this.Options.Ignore)
	}

	this.Output = &cli.Output{
		Console: make(chan string),
		Closed:  make(chan bool),
	}
}

func (this *Handler) Walk() error {
	return filepath.Walk(this.Options.Dir, func(fullPath string, info os.FileInfo, err error) error {

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
}

func (this *Handler) matchPath(fullPath string) {
	defer this.Group.Done()

	if this.Ignore != nil && this.Ignore.MatchString(fullPath) {
		return
	}

	if this.Pattern.MatchString(fullPath) {
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
