package walker

import (
	"io/fs"
	"path/filepath"
	"regexp"

	"github.com/nomad-software/findr/cli"
	"github.com/nomad-software/findr/file/path"
)

// Walker is the main file walker.
type Walker struct {
	Dir    string
	Glob   string
	Ignore *regexp.Regexp
	Regex  *regexp.Regexp
	Path   *path.Worker
}

// New creates a new file walker. This walker will find valid files and pass
// them to the path worker for matching against the supplied options.
// This walker will skip symlinks.
func New(opt *cli.Options) *Walker {
	var w = Walker{
		Dir:   opt.Dir,
		Glob:  opt.Glob,
		Regex: compile(opt.Regex, opt.Case),
		Path:  path.New(),
	}

	if opt.Ignore != "" {
		w.Ignore = compile(opt.Ignore, opt.Case)
	}

	return &w
}

// Walk starts walking through the directory specified in the options and starts
// passing valid files to the path worker.
func (w *Walker) Walk() error {
	go w.Path.StartQueue()

	err := filepath.WalkDir(w.Dir, func(fullPath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		info, err := entry.Info()

		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		w.Path.Queue <- path.Match{
			FullPath: fullPath,
			Ignore:   w.Ignore,
			Regex:    w.Regex,
			Glob:     w.Glob,
		}

		return nil
	})

	w.Path.StopQueue()

	return err
}

// compile checks that a regex pattern compiles and then returns it.
func compile(pattern string, observeCase bool) (regex *regexp.Regexp) {
	if observeCase {
		regex, _ = regexp.Compile(pattern)
	} else {
		regex, _ = regexp.Compile("(?i)" + pattern)
	}

	return regex
}
