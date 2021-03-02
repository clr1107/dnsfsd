package cmd

import (
	"fmt"
	"net/url"
	"path"

	"github.com/clr1107/dnsfsd/pkg/rules"
	"github.com/spf13/cobra"
)

var (
	downloadCmd = &cobra.Command{
		Use:   "download <url> <destination file name>",
		Short: "Download a third party rule file",
		Long:  `Download a text file containing rules from a remote network for use in this local dns server.`,
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

	patterns, err := rules.DownloadRuleFile(u.String(), args[1])

	if err != nil {
		return err
	}

	fmt.Printf("Downloaded '%v' to %v (%v patterns)\n", u.String(), path.Join("/etc/dnsfsd/rules", args[1]), patterns)
	return nil
}
