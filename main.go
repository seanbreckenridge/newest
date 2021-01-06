package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// helper to print an error and exit, if an error was received
func errExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

type NewestConfig struct {
	includeDirs bool
	dir         string
}

func parseFlags() (*NewestConfig, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: newest [-include-dirs] [dir]

Prints the newest file in a directory

If a positional argument is not provided, prints the newest file in this directory
If the current directory could not be determined, this fails

Optional arguments:
`)
		flag.PrintDefaults()
	}
	includeDirs := flag.Bool("include-dirs", false, "Include directories in addition to files")
	flag.Parse()
	var dir string
	switch flag.NArg() {
	case 0:
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		dir = cwd
	case 1:
		dir = flag.Arg(0)
	default:
		return nil, errors.New("Provided too many positional arguments")
	}
	return &NewestConfig{
		includeDirs: *includeDirs,
		dir:         dir,
	}, nil
}

func newestPath(inDir string, includeDirs bool) (os.FileInfo, error) {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		return nil, err
	}
	var recent os.FileInfo
	var recentTime time.Time
	for _, fi := range files {
		// if this is a directory and we're meant to ignore directories
		if fi.Mode().IsDir() && !includeDirs {
			continue
		}
		// if we haven't found any files that have matched yet
		if recent == nil {
			recent = fi
			recentTime = fi.ModTime()
			continue
		}
		// if the most recent time was before this files
		fiTime := fi.ModTime()
		if recentTime.Before(fiTime) {
			recent = fi
			recentTime = fiTime
		}
	}
	// didn't find any files with this pattern (i.e. -include-dirs)
	if recent == nil {
		return nil, fmt.Errorf("Could not find any matching files in %s", inDir)
	} else {
		return recent, nil
	}
}

func main() {
	conf, err := parseFlags()
	errExit(err)
	newest, err := newestPath(conf.dir, conf.includeDirs)
	errExit(err)
	fmt.Println(path.Join(conf.dir, newest.Name()))
}
