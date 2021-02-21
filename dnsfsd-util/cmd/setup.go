package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Create all the necessary directories & files",
		Long:  `Create all the necessary directories and files.`,
		RunE:  runSetupSubCommand,
	}
)

const defaultConfig string = `port: 53
forwards:
- '1.0.0.1:53'
- '1.1.1.1:53'
log: '/var/log/dnsfsd/log.txt'
verbose: false
`

func createDirectories(path string) error {
	return os.MkdirAll(path, 0755)
}

func runSetupSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	const etcDirectory string = "/etc/dnsfsd/rules"
	const logDirectory string = "/var/log/dnsfsd"
	const configDirectory string = "/etc/dnsfsd/config.yml"

	if err := createDirectories(etcDirectory); err != nil {
		return fmt.Errorf("could not create directory %v: %v", etcDirectory, err)
	}

	if err := createDirectories(logDirectory); err != nil {
		return fmt.Errorf("could not create logging directory %v: %v", logDirectory, err)
	}

	if err := ioutil.WriteFile(configDirectory, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("could not write default configuration to %v: %v", configDirectory, err)
	}

	fmt.Println("Setup finished!")
	return nil
}
