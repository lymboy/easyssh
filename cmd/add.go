package cmd

import (
	"easyssh/config"
	"easyssh/util"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	addGroup string
	addEnv   string
	addUser  string
	addIPs   string
	addPort  int
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server to configuration",
	Long: `Interactively add a new server to the EasySSH configuration file.

Supports adding multiple servers at once by entering multiple IP addresses.
Default values are used when you press Enter without input.`,
	Example: `  easyssh add
  easyssh add -e web -i "192.168.1.10 192.168.1.11"
  easyssh add -g prod -e api -u root -i "192.168.1.10,192.168.1.11"`,
	Run: func(cmd *cobra.Command, args []string) {
		runAddServer()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addGroup, "group", "g", "", "Server group")
	addCmd.Flags().StringVarP(&addEnv, "env", "e", "", "Environment/Service name")
	addCmd.Flags().StringVarP(&addUser, "user", "u", "", "SSH user")
	addCmd.Flags().StringVarP(&addIPs, "ips", "i", "", "IP addresses (comma or space separated)")
	addCmd.Flags().IntVarP(&addPort, "port", "p", 22, "SSH port")
}

func runAddServer() {
	fmt.Println()
	fmt.Println("  Add new server(s) to EasySSH")
	fmt.Println()

	// Get defaults
	defaultUser := util.GetCurrentUsername()
	defaultGroup := "Default"
	defaultEnv := "service"
	defaultPort := 22

	// Apply command line flags or defaults
	group := addGroup
	env := addEnv
	user := addUser
	ipsInput := addIPs
	port := addPort
	if port == 0 {
		port = defaultPort
	}

	// Interactive mode - prompt for all fields
	needInteractive := env == "" || ipsInput == ""

	if needInteractive {
		fmt.Println("  Press Enter to use default value in [brackets]")
		fmt.Println()

		// Prompt for Group
		if group == "" {
			group = util.PromptInput("Group", defaultGroup)
			if group == "" {
				group = defaultGroup
			}
		}

		// Prompt for Name (environment)
		if env == "" {
			env = util.PromptInput("Name", defaultEnv)
			if env == "" {
				env = defaultEnv
			}
		}

		// Prompt for User
		if user == "" {
			user = util.PromptInput("User", defaultUser)
			if user == "" {
				user = defaultUser
			}
		}

		// Prompt for IP addresses
		fmt.Println()
		fmt.Println("  Enter IP addresses (supports comma, space, Chinese comma separators)")
		if ipsInput == "" {
			ipsInput = util.PromptInput("IP addresses", "")
		}
	}

	ips := util.ParseIPs(ipsInput)
	if len(ips) == 0 {
		fmt.Println()
		fmt.Println("  ✗ No valid IP addresses found")
		os.Exit(1)
	}

	// Preview configuration
	fmt.Println()
	fmt.Println("  Configuration:")
	fmt.Printf("    Group: %s\n", group)
	fmt.Printf("    Name:  %s\n", env)
	fmt.Printf("    User:  %s\n", user)
	fmt.Printf("    Port:  %d\n", port)
	fmt.Println()

	// Generate server names
	var serversToAdd []config.Server
	if len(ips) == 1 {
		// Single IP: use env name directly
		serversToAdd = append(serversToAdd, config.Server{
			Group: group,
			Name:  env,
			Host:  ips[0],
			Port:  port,
			User:  user,
		})
		fmt.Printf("  Will add 1 server:\n")
		fmt.Printf("    %s/%s (%s)\n", group, env, ips[0])
	} else {
		// Multiple IPs: add suffix
		fmt.Printf("  Will add %d servers:\n", len(ips))
		for i, ip := range ips {
			serverName := fmt.Sprintf("%s-%d", env, i+1)
			serversToAdd = append(serversToAdd, config.Server{
				Group: group,
				Name:  serverName,
				Host:  ip,
				Port:  port,
				User:  user,
			})
			fmt.Printf("    %s/%s (%s)\n", group, serverName, ip)
		}
	}

	// Save to config file directly
	fmt.Println()
	configPath := getConfigPath()
	if err := appendServersToConfig(configPath, serversToAdd); err != nil {
		fmt.Printf("  ✗ Failed to save: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ✓ Saved %d server(s) to %s\n", len(ips), configPath)
	fmt.Println()
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	return util.GetHomeDir() + "/.easyssh/easy_config.yaml"
}

// appendServersToConfig appends new servers to the config file
func appendServersToConfig(configPath string, servers []config.Server) error {
	// Read existing config
	conf := config.GetConf()

	// Create a map of existing server names for deduplication
	existingKeys := make(map[string]bool)
	for _, s := range conf.GetServerList() {
		existingKeys[s.GetName()] = true
	}

	// Add new servers (skip duplicates)
	addedCount := 0
	for _, s := range servers {
		if !existingKeys[s.GetName()] {
			conf.ServerList = append(conf.ServerList, s)
			addedCount++
		}
	}

	if addedCount == 0 {
		return nil
	}

	// Write config back to file
	return writeConfigFile(configPath, conf)
}

// writeConfigFile writes the config to a YAML file
func writeConfigFile(path string, conf *config.Config) error {
	// Build YAML content
	var content string
	content += "ssh:\n"
	sshConfig := conf.GetSSHConfig()
	content += fmt.Sprintf("  key: %q\n", sshConfig.GetKey())
	content += fmt.Sprintf("  keep_alive: %s\n", strconv.FormatBool(sshConfig.KeepAlive))
	content += fmt.Sprintf("  keep_alive_interval: %q\n", sshConfig.KeepAliveInterval)
	content += fmt.Sprintf("  use_system_ssh: %s\n", strconv.FormatBool(sshConfig.UseSystemSSH))

	content += "\nserver:\n"
	for _, s := range conf.ServerList {
		content += fmt.Sprintf("  - group: %q\n", s.Group)
		content += fmt.Sprintf("    name: %q\n", s.Name)
		content += fmt.Sprintf("    host: %q\n", s.Host)
		if s.Port != 22 && s.Port != 0 {
			content += fmt.Sprintf("    port: %d\n", s.Port)
		}
		if s.User != "" {
			content += fmt.Sprintf("    user: %q\n", s.User)
		}
		if s.Password != "" {
			content += fmt.Sprintf("    password: %q\n", s.Password)
		}
		if s.Desc != "" {
			content += fmt.Sprintf("    desc: %q\n", s.Desc)
		}
	}

	// Ensure directory exists
	dir := util.GetHomeDir() + "/.easyssh"
	if !util.Exists(dir) {
		util.Mkdir(dir)
	}

	return os.WriteFile(path, []byte(content), 0644)
}
