package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "v1.1.0"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "version command.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
	},
}
