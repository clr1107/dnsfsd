package main

import (
	"regexp"

	"github.com/clr1107/dnsfsd/dnsfsd-util/cmd"
)

var (
	Patterns []*regexp.Regexp
)

func main() {
	cmd.ExecuteRoot()
}
