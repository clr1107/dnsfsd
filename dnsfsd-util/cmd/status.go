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
	s, err := active.Output()
	code := active.ProcessState.ExitCode()

	if code == 0 {
		return true, nil
	} else {
		if string(s) == "failed\n" {
			return false, nil
		}

		if err == nil {
			err = fmt.Errorf("unknown error whilst checking systemd")
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
