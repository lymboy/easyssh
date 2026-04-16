package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup SSH ControlMaster for connection reuse",
	Long: `Configure SSH ControlMaster to enable connection reuse across terminals.

This will add the following to your ~/.ssh/config:
  Host *
      ControlMaster auto
      ControlPath ~/.ssh/sockets/%r@%h-%p
      ControlPersist no
      ServerAliveInterval 60`,
	Example: `  easyssh setup`,
	Run: func(cmd *cobra.Command, args []string) {
		runSetup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup() {
	fmt.Println()
	fmt.Println("  EasySSH Setup - SSH ControlMaster Configuration")
	fmt.Println()

	// Check if ~/.ssh/config exists
	sshConfigPath := expandPath("~/.ssh/config")
	socketsDir := expandPath("~/.ssh/sockets")

	// Read existing config
	existingConfig := ""
	if data, err := os.ReadFile(sshConfigPath); err == nil {
		existingConfig = string(data)
	}

	// Check if already configured
	if containsControlMasterConfig(existingConfig) {
		fmt.Println("  ✓ SSH ControlMaster is already configured.")
		fmt.Println()
		fmt.Println("  Your ~/.ssh/config already contains ControlMaster settings.")
		return
	}

	// Show what will be added
	fmt.Println("  Will add the following to ~/.ssh/config:")
	fmt.Println()
	fmt.Println("    # EasySSH Connection Reuse")
	fmt.Println("    Host *")
	fmt.Println("        ControlMaster auto")
	fmt.Println("        ControlPath ~/.ssh/sockets/%r@%h-%p")
	fmt.Println("        ControlPersist no")
	fmt.Println("        ServerAliveInterval 60")
	fmt.Println()

	// Add configuration
	configToAdd := `
# EasySSH Connection Reuse
Host *
    ControlMaster auto
    ControlPath ~/.ssh/sockets/%r@%h-%p
    ControlPersist no
    ServerAliveInterval 60
`

	// Append to config
	newConfig := existingConfig + configToAdd

	// Ensure ~/.ssh directory exists
	sshDir := filepath.Dir(sshConfigPath)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		fmt.Printf("  ✗ Failed to create ~/.ssh directory: %s\n", err)
		os.Exit(1)
	}

	// Write config
	if err := os.WriteFile(sshConfigPath, []byte(newConfig), 0600); err != nil {
		fmt.Printf("  ✗ Failed to write config: %s\n", err)
		os.Exit(1)
	}

	// Create sockets directory
	if err := os.MkdirAll(socketsDir, 0700); err != nil {
		fmt.Printf("  ✗ Failed to create sockets directory: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("  ✓ SSH ControlMaster configured successfully!")
	fmt.Println()
	fmt.Println("  Created:")
	fmt.Printf("    • %s (updated)\n", sshConfigPath)
	fmt.Printf("    • %s (created)\n", socketsDir)
	fmt.Println()
	fmt.Println("  Next steps:")
	fmt.Println("    1. Enable in EasySSH config: use_system_ssh: true")
	fmt.Println("    2. Or run: easyssh config set use_system_ssh true")
	fmt.Println()
}

func containsControlMasterConfig(config string) bool {
	return contains(config, "ControlMaster auto") ||
		contains(config, "ControlPath")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func expandPath(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func restartSSHAgent() {
	// Optionally restart ssh-agent to apply changes
	exec.Command("pkill", "-HUP", "ssh-agent").Run()
}
