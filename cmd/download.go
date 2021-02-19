package cmd

import (
	"fmt"
	"net/url"

	"github.com/clr1107/dnsfsd/pkg/persistence"
	"github.com/spf13/cobra"
)

var (
	downloadCmd = &cobra.Command{
		Use:   "download <url> <destination file name>",
		Short: "Download a third party pattern",
		Long:  `Download a text file containing regular expressions as patterns from a remote network for use in this local dns server.`,
		RunE:  runDownloadSubCommand,
	}
)

func runDownloadSubCommand(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmd.Help()
	}

	u, err := url.ParseRequestURI(args[0])
	if err != nil {
		return err
	}

	patterns, err := persistence.DownloadPattern(u.String(), args[1])

	if err != nil {
		return err
	}

	fmt.Printf("Downloaded '%v' to %v (%v patterns)\n", u.String(), "/etc/dnsfsd/patterns/"+args[1], patterns)
	return nil
}
