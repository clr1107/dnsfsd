package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
)

import _ "embed"

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Systemd status for dnsfsd",
		Long:  `Output the systemd status for dnsfsd (if systemd is in operation).`,
		RunE:  runStatusSubCommand,
	}
)

func isActive() (bool, error) {
	active := exec.Command("systemctl", "is-active", "dnsfsd")
	err := active.Run()
	code := active.ProcessState.ExitCode()

	if code == 0 {
		return true, nil
	} else {
		if err == nil {
			return false, nil
		}

		return false, err
	}
}

func runStatusSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		return cmd.Help()
	}

	active, err := isActive()

	if err != nil {
		return fmt.Errorf("checking active: %v", err)
	}

	fmt.Printf("Systemd active: %v\n", active)

	return nil
}
