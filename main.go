package main

import (
	"regexp"

	"github.com/clr1107/dnsfsd/cmd"
)

var (
	Patterns []*regexp.Regexp
)

func main() {
	cmd.ExecuteRoot()
}
