package cmd

import (
	"fmt"

	"github.com/rajputpriyansh444-commits/flow/internal/cmn/config"
	"github.com/spf13/cobra"
)

func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the Flow version information",
		Long: `Print the current version number and build details of the Flow executable.

This command shows the version information that can be useful for troubleshooting,
reporting issues, or verifying that you're running the expected version of Flow.`,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(config.Version)
		},
	}
}
