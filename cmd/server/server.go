package server

import (
	"easyssh/config"
	"easyssh/ssh"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "ServerList commands",
	Long:  `ServerList commands`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			gotoSSH(args[0])
		} else {
			tip()
		}
	},
}

func tip() {
	fmt.Println("Usage: easyssh server [server-name]")
}

func init() {
	ServerCmd.AddCommand(ServerLsCmd)
}

func gotoSSH(str string) {
	server := config.GetConf().GetServer(str)
	err := ssh.New(server).RunTerminal(os.Stdout, os.Stderr)
	if err != nil {
		return
	}
}
