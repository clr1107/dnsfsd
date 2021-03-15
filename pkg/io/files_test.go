package io

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func writeToTempFile(data string) (*os.File, error) {
	f, err := ioutil.TempFile("", "dnsfsd_testing_pkg_io")
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(data)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func TestReadLines(t *testing.T) {
	const data string = "line one\nline two\nline three\nline four\n"
	const checkAgainst string = "line one\nline two\nline three\n"

	f, err := writeToTempFile(data)
	if err != nil {
		t.Fatalf("couldn't create and/or write to temp file in #TestReadLines: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	tempFile, _ := os.Open(f.Name())
	read := ReadFileLines(tempFile, 3)
	tempFile.Close()

	var merged strings.Builder
	for _, line := range read {
		merged.WriteString(string(line) + "\n")
	}

	mergedString := merged.String()
	if mergedString != checkAgainst {
		t.Fatalf("read string '%v' does not match data '%v'", mergedString, checkAgainst)
	}
}

func TestReadLinesReverse(t *testing.T) {
	const data string = "line one\nline two\nline three\nline four\n"
	const checkAgainst string = "line one\nline two\nline three\n"

	f, err := writeToTempFile(data)
	if err != nil {
		t.Fatalf("couldn't create and/or write to temp file in #TestReadLinesReverse: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	tempFile, _ := os.Open(f.Name())
	read := ReadFileLines(tempFile, 3)
	tempFile.Close()

	var merged strings.Builder
	for _, line := range read {
		merged.WriteString(string(line) + "\n")
	}

	mergedString := merged.String()
	if mergedString != checkAgainst {
		t.Fatalf("read string '%v' does not match data '%v'", mergedString, checkAgainst)
	}
}
