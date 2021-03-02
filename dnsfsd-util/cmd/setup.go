package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

import _ "embed"

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Create all the necessary directories & files",
		Long:  `Create all the necessary directories and files. This will overwrite any existing configuration with a default one.`,
		RunE:  runSetupSubCommand,
	}
	//go:embed static/config.yml
	defaultConfig []byte
	//go:embed static/dnsfsd.service
	serviceFile []byte
)

func createDirectories(path string) error {
	return os.MkdirAll(path, 0755)
}

func writeSystemd() error {
	return ioutil.WriteFile(
		"/etc/systemd/system/dnsfsd.service",
		serviceFile,
		0644,
	)
}

func exists(path string) bool {
	_, err := os.Stat(path);
	return err == nil
}

func runSetupSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	const etcDirectory string = "/etc/dnsfsd/rules"
	const logDirectory string = "/var/log/dnsfsd"
	const configPath string = "/etc/dnsfsd/config.yml"

	if exists(etcDirectory) {
		println(etcDirectory + " already exists")
	} else {
		if err := createDirectories(etcDirectory); err != nil {
			return fmt.Errorf("could not create directory %v: %v", etcDirectory, err)
		} else {
			println("created directory " + etcDirectory)
		}
	}

	if exists(logDirectory) {
		println(logDirectory + " already exists")
	} else {
		if err := createDirectories(logDirectory); err != nil {
			return fmt.Errorf("could not create logging directory %v: %v", logDirectory, err)
		} else {
			println("created directory " + logDirectory)
		}
	}

	if exists(configPath) {
		println(configPath + " already exists")
	} else {
		if err := ioutil.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return fmt.Errorf("could not write default configuration to %v: %v", configPath, err)
		} else {
			println("written default configuration to " + configPath)
		}
	}

	if exists("/etc/systemd/system/dnsfsd.service") {
		println("systemd service file already exists")
	} else {
		if err := writeSystemd(); err != nil {
			return fmt.Errorf("could not write service file /etc/systemd/system/dnsfsd.service")
		} else {
			println("written systemd service file")
		}
	}

	println()
	fmt.Printf("Configuration path: %v\n", configPath)
	fmt.Printf("Logs directory: %v\n", logDirectory)
	fmt.Println("Setup finished!")

	return nil
}
