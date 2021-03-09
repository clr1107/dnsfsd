package io

import (
	"bufio"
	"os"
)

const (
	CharNL byte = 10
	CharCR byte = 13
)

func ReadFileLines(file *os.File, count int64) [][]byte {
	var counter int64
	var read [][]byte

	if count > 0 {
		read = make([][]byte, 0, count)
	} else {
		read = make([][]byte, 0)
	}

	sc := bufio.NewScanner(file)

	for counter = 0; sc.Scan(); counter++ {
		if count >= 0 && counter >= count {
			break
		}

		read = append(read, sc.Bytes())
	}

	return read
}

func ReadFileLinesReverse(file *os.File, count int64) ([][]byte, error) {
	var err error
	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	var cursor int64
	var lines [][]byte
	size := stat.Size()
	line := make([]byte, 0)

	if count > 0 {
		lines = make([][]byte, 0, count)
	} else {
		lines = make([][]byte, 0)
	}

	for cursor = -1; int64(len(lines)) < count; cursor-- {
		if _, err := file.Seek(cursor, 2); err != nil {
			return nil, err
		}

		char := make([]byte, 1)
		if _, err := file.Read(char); err != nil {
			return nil, err
		}

		if char[0] == CharNL || char[0] == CharCR {
			if cursor != -1 {
				lines = append([][]byte{line}, lines...)
				line = make([]byte, 0)
			}
		} else {
			line = append(char, line...)
		}

		if cursor == -size {
			if len(line) != 0 {
				lines = append([][]byte{line}, lines...)
			}

			break
		}
	}

	return lines, nil
}
