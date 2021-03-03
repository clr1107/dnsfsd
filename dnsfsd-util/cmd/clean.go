package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

import _ "embed"

var (
	cleanCmd = &cobra.Command{
		Use:   "clean",
		Short: "Delete all logs",
		Long:  `Clean up disk space by deleting all logs.`,
		RunE:  runCleanSubCommand,
	}
)

func runCleanSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	const path string = "/var/log/dnsfsd/log.txt"
	return os.Remove(path)
}
