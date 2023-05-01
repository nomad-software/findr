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

// Options contain the command line options passed to the program.
type Options struct {
	Case   bool
	Dir    string
	Glob   string
	Help   bool
	Ignore string
	Regex  string
}

// ParseOptions parses the command line options.
func ParseOptions() *Options {
	var opt Options

	flag.BoolVar(&opt.Case, "case", false, "Use to switch to case sensitive pattern matching.")
	flag.StringVar(&opt.Dir, "dir", defaultDirectory, "The directory to traverse.")
	flag.StringVar(&opt.Glob, "glob", defaultGlob, "The glob file pattern to match.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.StringVar(&opt.Regex, "regex", defaultRegex, "A regex to match files against.")
	flag.Parse()

	opt.Dir, _ = homedir.Expand(opt.Dir)

	return &opt
}

// Valid checks command line options are valid.
func (opt *Options) Valid() bool {

	if err := compile(opt.Regex, opt.Case); err != nil {
		fmt.Printf("find pattern: %s", err.Error())
		return false
	}

	if err := compile(opt.Ignore, opt.Case); err != nil {
		fmt.Printf("ignore pattern: %s", err.Error())
		return false
	}

	return true
}

// compile checks that a regex pattern compiles.
func compile(pattern string, observeCase bool) (err error) {
	if observeCase {
		_, err = regexp.Compile(pattern)
	} else {
		_, err = regexp.Compile("(?i)" + pattern)
	}

	return err
}
