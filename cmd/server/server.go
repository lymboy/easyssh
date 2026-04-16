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
	Short: "Connect to a server",
	Long:  `Connect to a server by index number or name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			gotoSSH(args[0])
		} else {
			printUsage()
		}
	},
}

func printUsage() {
	fmt.Println()
	fmt.Println("Usage: easyssh server <index|name>")
	fmt.Println("       easyssh server ls")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  easyssh server 0        # Connect to first server")
	fmt.Println("  easyssh server web-prod  # Connect by name")
	fmt.Println("  easyssh server ls        # List all servers")
	fmt.Println()
}

func init() {
	ServerCmd.AddCommand(ServerLsCmd)
}

func gotoSSH(str string) {
	conf := config.GetConf()
	server := conf.GetServer(str)

	if server == nil {
		printServerNotFoundError(str, conf)
		os.Exit(1)
		return
	}

	err := ssh.New(server).RunTerminal(os.Stdout, os.Stderr)
	if err != nil {
		printConnectionError(server, err)
		os.Exit(1)
	}
}

func printServerNotFoundError(input string, conf *config.Config) {
	fmt.Println()
	printErrorBox(fmt.Sprintf("找不到服务器: '%s'", input))
	fmt.Println()

	// Suggest similar names
	similar := conf.GetSimilarServerNames(input, 2, 3)
	if len(similar) > 0 {
		fmt.Println("  您是否想连接以下服务器？")
		for _, name := range similar {
			svr := conf.GetServerByName(name)
			if svr != nil {
				fmt.Printf("  • %s (%s@%s)\n", colorYellow(name), svr.GetUser(), svr.GetHost())
			}
		}
		fmt.Println()
	}

	fmt.Println(colorDim("使用 'easyssh server ls' 查看完整列表"))
	printErrorBoxEnd()
	fmt.Println()
}

func printConnectionError(server *config.Server, err error) {
	fmt.Println()
	printErrorBox(fmt.Sprintf("连接失败: %s@%s:%d", server.GetUser(), server.GetHost(), server.GetPort()))
	fmt.Println()
	fmt.Printf("  错误: %s\n", colorRed(err.Error()))
	fmt.Println()
	fmt.Println(colorDim("请检查:"))
	fmt.Println(colorDim("  • 服务器地址和端口是否正确"))
	fmt.Println(colorDim("  • 用户名和认证方式是否正确"))
	fmt.Println(colorDim("  • 网络是否可达"))
	printErrorBoxEnd()
	fmt.Println()
}

func printErrorBox(msg string) {
	fmt.Print(colorCyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
	fmt.Printf("  %s %s\n", colorRed("✗"), msg)
	fmt.Print(colorCyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
}

func printErrorBoxEnd() {
	fmt.Print(colorCyan("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
}

func colorRed(s string) string {
	return fmt.Sprintf("\033[1;31m%s\033[0m", s)
}

func colorYellow(s string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", s)
}

func colorCyan(s string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", s)
}

func colorDim(s string) string {
	return fmt.Sprintf("\033[2m%s\033[0m", s)
}
