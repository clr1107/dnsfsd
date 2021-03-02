package cmd

import (
	"fmt"
	"strings"

	"github.com/clr1107/dnsfsd/pkg/rules"
	"github.com/spf13/cobra"
)

var (
	remove bool

	patternsCmd = &cobra.Command{
		Use:   "rules",
		Short: "List rules. If there are a large amount of rules this could take a long time!",
		Long:  `List all the rules currently being matched. If there are a large amount of rules this could take a long time!`,
		RunE:  runPatternsSubCommand,
	}
)

func loadRules() (*[]rules.RuleFile, error) {
	return rules.LoadAllRuleFiles("/etc/dnsfsd/rules")
}

func runPatternsSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	files, err := loadRules()

	if err != nil {
		return err
	}

	numOfRules := 0

	for _, v := range *files {
		numOfRules += len(*v.Rules)
	}

	header := fmt.Sprintf("Matching %v rules", numOfRules)
	println(header)
	println(strings.Repeat("=", len(header)))

	for _, i := range *files {
		println()

		header := fmt.Sprintf("File '%v' (%v)", i.Path, len(*i.Rules))
		println(header)
		println(strings.Repeat("-", len(header)))

		for k, j := range *i.Rules {
			fmt.Printf("%v)    %v\n", k+1, j)
		}
	}

	return nil
}
