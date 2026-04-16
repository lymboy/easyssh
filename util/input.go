package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// PromptInput prompts user for input with a default value
// Press Enter to use default value (conda style)
func PromptInput(prompt string, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	// Display prompt with default in brackets (conda style)
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}
	return input
}

// ParseIPs parses multiple IP addresses from input
// Supports: space, comma (both half-width and full-width), newline, tab separators
func ParseIPs(input string) []string {
	// Replace all separators with space
	input = strings.ReplaceAll(input, "，", " ")  // full-width comma
	input = strings.ReplaceAll(input, ",", " ")  // half-width comma
	input = strings.ReplaceAll(input, "\t", " ") // tab
	input = strings.ReplaceAll(input, "\n", " ") // newline

	// Split by space
	parts := strings.Fields(input)

	// Validate and deduplicate
	seen := make(map[string]bool)
	var validIPs []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || seen[part] {
			continue
		}

		// Accept IP addresses and hostnames
		if isValidIPOrHost(part) {
			validIPs = append(validIPs, part)
			seen[part] = true
		}
	}

	return validIPs
}

// isValidIPOrHost checks if the string is a valid IP address or hostname
func isValidIPOrHost(s string) bool {
	// Check if it's a valid IP address
	if ip := net.ParseIP(s); ip != nil {
		return true
	}

	// Check if it's a valid hostname (simple check)
	// Hostname regex: alphanumeric and hyphens, with dots
	hostnameRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	return hostnameRegex.MatchString(s) && len(s) <= 253
}

// PromptYesNo prompts for yes/no confirmation (conda style: [y/N] or [Y/n])
func PromptYesNo(prompt string, defaultYes bool) bool {
	reader :=bufio.NewReader(os.Stdin)

	var hint string
	if defaultYes {
		hint = "Y/n"
	} else {
		hint = "y/N"
	}

	fmt.Printf("%s [%s]: ", prompt, hint)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// GetCurrentUsername returns the current system username
func GetCurrentUsername() string {
	user, err := os.UserHomeDir()
	if err != nil {
		return os.Getenv("USER")
	}
	parts := strings.Split(user, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return os.Getenv("USER")
}
