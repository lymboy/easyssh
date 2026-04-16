package cmd

import (
	"easyssh/config"
	"easyssh/util"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <server>",
	Short:   "Remove a server from configuration",
	Long:    `Remove a server from the EasySSH configuration file.`,
	Example: `  easyssh remove 0        # Remove by index
  easyssh remove web-1    # Remove by name
  easyssh remove uat      # Remove by name`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runRemoveServer(args[0])
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemoveServer(input string) {
	conf := config.GetConf()
	servers := conf.GetServerList()

	if len(servers) == 0 {
		fmt.Println("  No servers configured.")
		return
	}

	// Check if input is a number (index)
	var target *config.Server
	var targetIndex int = -1

	if util.IsDigit(input) {
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("  ✗ Invalid index: %s\n", input)
			os.Exit(1)
		}
		if index < 0 || index >= len(servers) {
			fmt.Printf("  ✗ Index %d out of range (0-%d)\n", index, len(servers)-1)
			os.Exit(1)
		}
		targetIndex = index
		target = &servers[index]
	} else {
		// Find by name
		for i, svr := range servers {
			if svr.GetName() == input {
				targetIndex = i
				target = &servers[i]
				break
			}
		}
	}

	if target == nil {
		fmt.Printf("  ✗ Server '%s' not found\n", input)

		// Suggest similar names
		similar := conf.GetSimilarServerNames(input, 2, 3)
		if len(similar) > 0 {
			fmt.Println()
			fmt.Println("  Did you mean:")
			for _, name := range similar {
				fmt.Printf("    • %s\n", name)
			}
		}
		os.Exit(1)
	}

	// Show server to remove
	fmt.Println()
	fmt.Printf("  Removing server:\n")
	fmt.Printf("    %s/%s (%s@%s)\n", target.GetGroup(), target.GetName(), target.GetUser(), target.GetHost())
	fmt.Println()

	// Remove from slice
	newServers := make([]config.Server, 0, len(servers)-1)
	for i, svr := range servers {
		if i != targetIndex {
			newServers = append(newServers, svr)
		}
	}
	conf.ServerList = newServers

	// Save to config file
	configPath := getConfigPath()
	if err := writeConfigFile(configPath, conf); err != nil {
		fmt.Printf("  ✗ Failed to save: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ✓ Removed '%s' from configuration\n", target.GetName())
	fmt.Println()
}

func getConfigPathForRemove() string {
	return util.GetHomeDir() + "/.easyssh/easy_config.yaml"
}
