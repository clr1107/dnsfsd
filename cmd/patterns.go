package cmd

import (
	"fmt"
	"strings"

	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/spf13/cobra"
)

var (
	remove bool

	patternsCmd = &cobra.Command{
		Use:   "patterns",
		Short: "List patterns",
		Long:  `List all the patterns currently being matched.`,
		RunE:  runPatternsSubCommand,
	}
)

func loadPatterns() ([]*persistence.PatternFile, error) {
	return persistence.LoadAllPatternFiles("/etc/dnsfsd/patterns")
}

func runPatternsSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	files, err := loadPatterns()

	if err != nil {
		return err
	}

	numOfPatterns := 0

	for _, v := range files {
		numOfPatterns += len(v.Patterns)
	}

	header := fmt.Sprintf("Matching %v patterns", numOfPatterns)
	println(header)
	println(strings.Repeat("=", len(header)))

	for _, i := range files {
		println()

		header := fmt.Sprintf("File '%v'", i.Path)
		println(header)
		println(strings.Repeat("-", len(header)))

		for k, j := range i.Patterns {
			fmt.Printf("%v)  '%v'\n", k+1, j)
		}

	}

	return nil
}
