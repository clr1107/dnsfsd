package persistence

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
		return fmt.Errorf("could not open pattern file '%v'", p.Path)
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	patterns := make([]*regexp.Regexp, 0)

	for scanner.Scan() {
		text := scanner.Text()
		pattern, err := regexp.Compile(text)

		if err != nil {
			return fmt.Errorf("could not parse regular expression '%v'", text)
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

func CollectAllPatterns(files []*PatternFile) []*regexp.Regexp {
	l := make([]*regexp.Regexp, 0, len(files))

	for _, v := range files {
		if v.Loaded {
			l = append(l, v.Patterns...)
		}
	}

	return l
}

func DownloadPattern(url string, filename string) (int, error) {
	resp, err := http.Get(url)

	if err != nil {
		return 0, fmt.Errorf("could not download from url: '%v'", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("HTTP request to '%v' gave a %v status code", url, resp.StatusCode)
	}

	const directory string = "/etc/dnsfsd/patterns"
	if err := os.MkdirAll(directory, os.FileMode(0755)); err != nil {
		return 0, err
	}

	filepath := directory + string(os.PathSeparator) + filename
	out, err := os.Create(filepath)

	if err != nil {
		return 0, fmt.Errorf("could not create file '%v'", filepath)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return 0, err
	}

	patternFile := &PatternFile{filepath, false, nil}

	if err := patternFile.Load(); err != nil {
		return 0, err
	} else {
		return len(patternFile.Patterns), nil
	}
}
