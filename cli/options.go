package cli

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/mitchellh/go-homedir"
)

const (
	defaultDirectory = "."
	defaultRegex     = ".*"
	defaultGlob      = "*"
)

// Options contains CLI arguments passed to the program.
type Options struct {
	Case   bool
	Dir    string
	Glob   string
	Help   bool
	Ignore string
	Regex  string
}

// ParseOptions parses the command line options and returns a struct filled with
// the relevant options.
func ParseOptions() Options {
	var opt Options

	flag.BoolVar(&opt.Case, "case", false, "Use to switch to case sensitive pattern matching.")
	flag.StringVar(&opt.Dir, "dir", defaultDirectory, "The directory to traverse.")
	flag.StringVar(&opt.Glob, "glob", defaultGlob, "The glob file pattern to match.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.StringVar(&opt.Regex, "regex", defaultRegex, "A regex to match files against.")
	flag.Parse()

	opt.Dir, _ = homedir.Expand(opt.Dir)

	return opt
}

// Valid checks command line options are valid.
func (opt *Options) Valid() bool {

	err := opt.compiles(opt.Regex)
	if err != nil {
		fmt.Printf("find pattern: %s", err.Error())
		return false
	}

	err = opt.compiles(opt.Ignore)
	if err != nil {
		fmt.Printf("ignore pattern: %s", err.Error())
		return false
	}

	return true
}

// PrintUsage prints the usage of the program.
func (opt *Options) PrintUsage() {
	var banner string = `_____ _           _
|  ___(_)_ __   __| |_ __
| |_  | | '_ \ / _' | '__|
|  _| | | | | | (_| | |
|_|   |_|_| |_|\__,_|_|

`
	fmt.Println(banner)
	flag.Usage()
}

// Check that a regex pattern compiles.
func (opt *Options) compiles(pattern string) (err error) {
	if opt.Case {
		_, err = regexp.Compile(pattern)
	} else {
		_, err = regexp.Compile("(?i)" + pattern)
	}

	return err
}
