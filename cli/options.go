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
	Dir     string
	Pattern string
	Ignore  string
	Case    bool
	Help    bool
}

func (this *Options) Parse() {
	flag.StringVar(&this.Dir, "dir", ".", "The directory to traverse.")
	flag.StringVar(&this.Pattern, "pattern", "", "The regex to match files against.")
	flag.StringVar(&this.Ignore, "ignore", "", "A regex to ignore files or directories.")
	flag.BoolVar(&this.Case, "case", false, "Use to switch to case sensitive matching.")
	flag.BoolVar(&this.Help, "help", false, "Show help.")
	flag.Parse()

	dir, _ := homedir.Expand(this.Dir)
	this.Dir = dir
}

func (this *Options) Valid() bool {

	if this.Pattern == "" {
		fmt.Fprintln(os.Stderr, color.RedString("A pattern must be defined."))
		return false
	}

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
	options := color.CyanString("finding:     ")
	options += color.GreenString("%s\n", this.Pattern)
	options += color.CyanString("starting in: ")
	options += color.GreenString("%s\n", this.Dir)

	if this.Ignore != "" {
		options += color.CyanString("ignoring:    ")
		options += color.GreenString("%s\n", this.Ignore)
	}

	fmt.Print(options)
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
