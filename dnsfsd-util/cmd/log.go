package cmd

import (
	"fmt"
	"github.com/clr1107/dnsfsd/pkg/data/config"
	"github.com/clr1107/dnsfsd/pkg/io"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

import _ "embed"

var (
	logCmd = &cobra.Command{
		Use:   "log [(-)length]",
		Short: "Output the log",
		Long:  `Output the log file, if it exists.`,
		RunE:  runLogSubCommand,
	}
)

func runLogSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return cmd.Help()
	}

	var length int64 = -1
	var reverse bool = false

	if len(args) == 1 {
		signedLength, err := strconv.ParseInt(args[0], 10, 64)

		if err != nil {
			return err
		}

		if signedLength < 0 {
			reverse = true
			length = -1 * signedLength
		} else {
			length = signedLength
		}
	}

	if err := config.InitConfig(); err != nil {
		return err
	}

	path := viper.GetString("log.path")
	fp, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("could not open log file: %v", err)
	}

	var lines [][]byte

	if reverse {
		lines, err = io.ReadFileLinesReverse(fp, length)

		if err != nil {
			return err
		}
	} else {
		lines = io.ReadFileLines(fp, length)
	}

	for _, line := range lines {
		println(string(line))
	}

	return nil
}
