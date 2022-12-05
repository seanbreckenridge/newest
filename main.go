package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"
)

type NewestConfig struct {
	includeDirs  bool
	ignoreHidden bool
	dir          string
}

func parseFlags() (*NewestConfig, error) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: newest [OPTIONS] [DIR]

Prints the newest file in a directory

If a positional argument is not provided, prints the newest file in this directory
If the current directory could not be determined, this fails

Optional arguments:`)
		flag.PrintDefaults()
	}
	includeDirs := flag.Bool("include-dirs", false, "Include directories in addition to files")
	ignoreHidden := flag.Bool("ignore-hidden", false, "Ignore hidden files")
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
		return nil, errors.New("newest: Provided too many positional arguments")
	}
	if *ignoreHidden && runtime.GOOS == "windows" {
		return nil, errors.New("newest: Not able to check if files are hidden on windows")
	}
	return &NewestConfig{
		includeDirs:  *includeDirs,
		ignoreHidden: *ignoreHidden,
		dir:          dir,
	}, nil
}

func newestPath(conf *NewestConfig) (os.FileInfo, error) {
	files, err := ioutil.ReadDir(conf.dir)
	if err != nil {
		return nil, err
	}
	var recent os.FileInfo
	var recentTime time.Time
	for _, fi := range files {
		// if this is a directory and we're meant to ignore directories
		if !conf.includeDirs && fi.Mode().IsDir() {
			continue
		}
		// if this is a hidden file and we're meant to ignore hidden files
		if conf.ignoreHidden && fi.Name()[0:1] == "." {
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
		return nil, fmt.Errorf("newest: Could not find any matching files in %s", conf.dir)
	} else {
		return recent, nil
	}
}

// wrapper for 'main' code, to return single err to main
func newest() (string, error) {
	conf, err := parseFlags()
	if err != nil {
		return "", err
	}
	newest, err := newestPath(conf)
	if err != nil {
		return "", err
	}
	return path.Join(conf.dir, newest.Name()), nil
}

func main() {
	result, err := newest()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println(result)
}
