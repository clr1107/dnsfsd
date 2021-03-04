package cmd

import (
	"fmt"
	"github.com/clr1107/dnsfsd/pkg/data/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

import _ "embed"

var (
	logCmd = &cobra.Command{
		Use:   "log",
		Short: "Output the log",
		Long:  `Output the log file, if it exists.`,
		RunE:  runLogSubCommand,
	}
)

func runLogSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	if err := config.InitConfig(); err != nil {
		return err
	}

	path := viper.GetString("log.path")
	b, err := ioutil.ReadFile(path)

	if err != nil {
		return fmt.Errorf("could not read log: %v", err)
	}

	print(string(b))
	return nil
}
