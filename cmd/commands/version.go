package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd(gitCommit, buildTime string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the application version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "Git Commit: %s\n", gitCommit)
			fmt.Fprintf(cmd.OutOrStdout(), "Build Time: %s\n", buildTime)
		},
	}
}
