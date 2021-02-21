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
		Long:  `Create all the necessary directories and files. This will overwrite any existing configuration with a default one.`,
		RunE:  runSetupSubCommand,
	}
)

// lol I'm lazy.
const defaultConfig string = `server:
    port: 53
    parallel_match: false
log:
    path: '/var/log/dnsfsd/log.txt'
    verbose: false
dns:
    cache: 86400
    forwards:
    - '1.0.0.1:53'
    - '1.1.1.1:53'
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
	const configPath string = "/etc/dnsfsd/config.yml"

	if err := createDirectories(etcDirectory); err != nil {
		return fmt.Errorf("could not create directory %v: %v", etcDirectory, err)
	}

	if err := createDirectories(logDirectory); err != nil {
		return fmt.Errorf("could not create logging directory %v: %v", logDirectory, err)
	}

	if err := ioutil.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("could not write default configuration to %v: %v", configPath, err)
	}

	fmt.Println("Setup finished!")
	fmt.Printf("Configuration path: %v\n", configPath)
	fmt.Printf("Logs directory: %v\n", logDirectory)

	return nil
}
