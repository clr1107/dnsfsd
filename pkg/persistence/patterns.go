package persistence

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type PatternFile struct {
	Path     string
	Loaded   bool
	Patterns []*regexp.Regexp
}

func (p *PatternFile) Load() error {
	f, err := os.Open(p.Path)

	if err != nil {
		return err
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	patterns := make([]*regexp.Regexp, 0)

	for scanner.Scan() {
		text := scanner.Text()
		pattern, err := regexp.Compile(text)

		if err != nil {
			return err
		}

		patterns = append(patterns, pattern)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	p.Patterns = patterns
	p.Loaded = true

	return nil
}

func AllPatternsFiles(path string) ([]*PatternFile, error) {
	files, err := ioutil.ReadDir(path)
	paths := make([]*PatternFile, 0)

	if err != nil {
		return nil, fmt.Errorf("could not read all pattern files in directory '%v'", path)
	}

	for _, v := range files {
		if !v.IsDir() {
			paths = append(paths, &PatternFile{path + string(os.PathSeparator) + v.Name(), false, nil})
		}
	}

	return paths, nil
}

func LoadAllPatternFiles(path string) ([]*PatternFile, error) {
	files, err := AllPatternsFiles(path)
	successes := make([]*PatternFile, 0, len(files))

	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if err := v.Load(); err == nil {
			successes = append(successes, v)
		} else {
			return nil, err
		}
	}

	return successes, nil
}
