package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "version command.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: v1.0.0")
	},
}
