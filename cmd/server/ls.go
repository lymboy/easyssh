package server

import (
	"easyssh/config"

	"github.com/spf13/cobra"
)

var ServerLsCmd = &cobra.Command{
	Use:          "ls",
	Short:        "ls server list",
	Example:      "easyssh server ls",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		print()
	},
}

func print() {
	config.GetConf().Print()
}
