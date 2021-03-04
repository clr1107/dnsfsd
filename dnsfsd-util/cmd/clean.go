package cmd

import (
	"github.com/clr1107/dnsfsd/pkg/data/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	if err := config.InitConfig(); err != nil {
		return err
	}

	path := viper.GetString("log.path")
	return os.Remove(path)
}
