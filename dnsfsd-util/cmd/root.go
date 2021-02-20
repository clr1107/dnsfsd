package cmd

import "github.com/spf13/cobra"

var (
	rootCmd = &cobra.Command{
		Use:          "dnsfs",
		Short:        "DNS Filtering Sinkhole",
		Long:         `A DNS server, designed to be run locally, that filters requests based on regular expressions (pattern matching) and either forwards the requests to another DNS server or simply filters them and ignores them -- sinkholes them.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func init() {
	rootCmd.AddCommand(patternsCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(setupCmd)
}

func ExecuteRoot() error {
	return rootCmd.Execute()
}
