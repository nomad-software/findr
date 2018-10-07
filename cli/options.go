package cli

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

const (
	DEFAULT_DIRECTORY = "."
	DEFAULT_REGEX     = ".*"
	DEFAULT_GLOB      = "*"
)

type Options struct {
	Case   bool
	Dir    string
	Glob   string
	Help   bool
	Ignore string
	Regex  string
}

func ParseOptions() Options {
	var opt Options

	flag.BoolVar(&opt.Case, "case", false, "Use to switch to case sensitive pattern matching.")
	flag.StringVar(&opt.Dir, "dir", DEFAULT_DIRECTORY, "The directory to traverse.")
	flag.StringVar(&opt.Glob, "glob", DEFAULT_GLOB, "The glob file pattern to match.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.StringVar(&opt.Regex, "regex", DEFAULT_REGEX, "A regex to match files against.")
	flag.Parse()

	opt.Dir, _ = homedir.Expand(opt.Dir)

	return opt
}

func (this *Options) Valid() bool {

	err := this.compiles(this.Regex)
	if err != nil {
		fmt.Fprintln(Stderr, color.RedString("find pattern: %s", err.Error()))
		return false
	}

	err = this.compiles(this.Ignore)
	if err != nil {
		fmt.Fprintln(Stderr, color.RedString("ignore pattern: %s", err.Error()))
		return false
	}

	return true
}

func (this *Options) Usage() {
	var banner string = `_____ _           _
|  ___(_)_ __   __| |_ __
| |_  | | '_ \ / _' | '__|
|  _| | | | | | (_| | |
|_|   |_|_| |_|\__,_|_|

`
	color.Cyan(banner)
	flag.Usage()
}

func (this *Options) compiles(pattern string) (err error) {
	if this.Case {
		_, err = regexp.Compile(pattern)
	} else {
		_, err = regexp.Compile("(?i)" + pattern)
	}

	return err
}
