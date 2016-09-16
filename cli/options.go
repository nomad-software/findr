package cli

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

type Options struct {
	Case    bool
	Dir     string
	File    string
	Help    bool
	Ignore  string
	Pattern string
}

func ParseOptions() Options {
	var opt Options

	flag.BoolVar(&opt.Case, "case", false, "Use to switch to case sensitive pattern matching.")
	flag.StringVar(&opt.Dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&opt.File, "file", "*", "The glob file pattern to match.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.StringVar(&opt.Pattern, "pattern", ".*", "A regex to match files against.")
	flag.Parse()

	opt.Dir, _ = homedir.Expand(opt.Dir)

	return opt
}

func (this *Options) Valid() bool {

	err := this.compiles(this.Pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("find pattern: %s", err.Error()))
		return false
	}

	err = this.compiles(this.Ignore)
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString("ignore pattern: %s", err.Error()))
		return false
	}

	return true
}

func (this *Options) Echo() {

	var output string

	if this.Pattern != ".*" {
		output += color.CyanString("finding:     ")
		output += color.GreenString("%s\n", this.Pattern)
	}

	if this.File != "*" {
		output += color.CyanString("files:       ")
		output += color.GreenString("%s\n", this.File)
	}

	output += color.CyanString("starting in: ")
	output += color.GreenString("%s\n", this.Dir)

	if this.Ignore != "" {
		output += color.CyanString("ignoring:    ")
		output += color.GreenString("%s\n", this.Ignore)
	}

	fmt.Print(output)
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
