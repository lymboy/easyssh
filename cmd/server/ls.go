package server

import (
	"easyssh/config"

	"github.com/spf13/cobra"
)

var ServerLsCmd = &cobra.Command{
	Use:          "ls",
	Short:        "List all servers",
	Example:      "easyssh server ls",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		config.GetConf().Print()
	},
}
