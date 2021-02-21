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
		Use:   "rules",
		Short: "List rules",
		Long:  `List all the rules currently being matched.`,
		RunE:  runPatternsSubCommand,
	}
)

func loadRules() (*[]persistence.RuleFile, error) {
	return persistence.LoadAllRuleFiles("/etc/dnsfsd/rules")
}

func runPatternsSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	files, err := loadRules()

	if err != nil {
		return err
	}

	numOfPatterns := 0

	for _, v := range *files {
		numOfPatterns += len(*v.Rules)
	}

	header := fmt.Sprintf("Matching %v rules", numOfPatterns)
	println(header)
	println(strings.Repeat("=", len(header)))

	for _, i := range *files {
		println()

		header := fmt.Sprintf("File '%v'", i.Path)
		println(header)
		println(strings.Repeat("-", len(header)))

		for k, j := range *i.Rules {
			fmt.Printf("%v)    %v\n", k+1, j)
		}

	}

	return nil
}
