package cmd

import (
	"fmt"
	"github.com/clr1107/dnsfsd/pkg/rules"
	"github.com/spf13/cobra"
	"strings"
)

var (
	digCmd = &cobra.Command{
		Use:   "dig <domain>",
		Short: "Run a DNS query through the server",
		Long:  `Run a DNS query (A) through the server and find out what it would do. As a form of testing.`,
		RunE:  runDigSubCommand,
	}
)

func runDigSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return cmd.Help()
	}

	domain := strings.ToLower(args[0])
	files, err := loadRules()

	if err != nil {
		return err
	}

	ruleset := rules.CollectAllRules(files)

	println("; Test DNS Ruleset")
	fmt.Printf("; checking against %v rules", ruleset.Size())
	fmt.Printf("; (A) %v\n;\n", domain)

	var test bool

	delta := timeIt(func() {
		test = ruleset.Test(domain)
	}).Milliseconds()

	if test {
		fmt.Printf("; ruleset indicates SINK on domain %v\n", domain)
	} else {
		fmt.Printf("; ruleset indicates to FORWARD domain %v to DNS server(s)\n", domain)
	}

	fmt.Printf(";\n; test took %v ms\n", delta)
	fmt.Println("; remember, this tool does not use any caching!")

	if delta > 10 { // arbitrary number tbh.
		fmt.Println("; this is an impairing response time, perhaps you have a large or overly complex ruleset?")
	} else if delta <= 5 { // once again, arbitrary lol.
		fmt.Println("; this is an excellent response time, it will not impair noticeable performance at all")
	} else {
		fmt.Printf("; this is indicative of good performance. this is an adequately sized ruleset.")
	}

	return nil
}
