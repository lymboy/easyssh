package config

import (
	"easyssh/util"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var once sync.Once
var conf = new(Config)

func GetConf() *Config {
	once.Do(func() {
		var err error = nil
		conf, err = readConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
	return conf
}

type Config struct {
	SSHConfig  SSHConfig `json:"ssh" yaml:"ssh" mapstructure:"ssh"`
	ServerList []Server  `json:"server" yaml:"server" mapstructure:"server"`
}

func (c Config) GetSSHConfig() *SSHConfig {
	if util.IsZero(c.SSHConfig) {
		return new(SSHConfig)
	}
	return &c.SSHConfig
}

func (c Config) GetServerList() []Server {
	if len(c.ServerList) == 0 {
		return make([]Server, 0)
	}
	// 按照 group, name 属性排序
	sort.Slice(c.ServerList, func(i, j int) bool {
		c.ServerList[i].GetGroup()

		groupA := strings.ToLower(c.ServerList[i].GetGroup())
		groupB := strings.ToLower(c.ServerList[j].GetGroup())

		if strings.EqualFold(groupA, groupB) {
			nameA := strings.ToLower(c.ServerList[i].GetName())
			nameB := strings.ToLower(c.ServerList[j].GetName())
			return nameA < nameB
		}
		return groupA < groupB
	})
	return c.ServerList
}

func (c Config) GetServer(str string) *Server {
	if util.IsDigit(str) {
		return c.GetServerByIndex(cast.ToInt(str))
	}
	return c.GetServerByName(str)
}

func (c Config) GetServerByName(name string) *Server {
	if len(name) == 0 {
		return nil
	}
	for _, server := range c.GetServerList() {
		if strings.EqualFold(server.GetName(), name) {
			return &server
		}
	}
	return nil
}

func (c Config) GetServerByIndex(index int) *Server {
	if index < 0 {
		return nil
	}
	if index < len(c.GetServerList()) {
		return &c.GetServerList()[index]
	}
	return nil
}

func (c Config) GetServerForMap() map[string]*Server {
	serverMap := make(map[string]*Server)
	for _, server := range c.GetServerList() {
		serverMap[server.GetName()] = &server
	}
	return serverMap
}

// GetAllServerNames returns all server names for similarity matching
func (c Config) GetAllServerNames() []string {
	names := make([]string, 0, len(c.GetServerList()))
	for _, server := range c.GetServerList() {
		names = append(names, server.GetName())
	}
	return names
}

// GetSimilarServerNames returns server names similar to the given name
func (c Config) GetSimilarServerNames(name string, maxDistance int, limit int) []string {
	if limit <= 0 {
		return nil
	}
	allNames := c.GetAllServerNames()
	similar := util.SimilarNames(name, allNames, maxDistance)
	if len(similar) > limit {
		similar = similar[:limit]
	}
	return similar
}

func (c Config) ContainsServer(str string) bool {
	if util.IsDigit(str) {
		return nil == c.GetServerByIndex(cast.ToInt(str))
	}
	return nil == c.GetServerByName(str)
}

func (c Config) PrintServer() {
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>  ServerList  <<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println()
	for i, server := range c.GetServerList() {
		fmt.Printf("\t[%d]\t%s\t\t%s\n", i, server.GetName(), server)
	}
}

func (c Config) PrintServerV2() {
	servers := c.GetServerList()

	// Calculate dynamic column widths using display width
	// NAME: increase minWidth for Chinese service names (e.g., "企业账号中心" = 12 display width)
	maxName := 15  // Minimum width for NAME column
	maxGroup := 8  // Minimum width for GROUP column
	maxHost := 12  // Minimum width for HOST column
	maxUser := 6   // Minimum width for USER column

	const maxNameWidth = 30
	const maxGroupWidth = 15
	const maxHostWidth = 35
	const maxUserWidth = 15

	for _, svr := range servers {
		maxName = maxInt(maxName, minInt(getDisplayWidth(svr.GetName()), maxNameWidth))
		maxGroup = maxInt(maxGroup, minInt(getDisplayWidth(svr.GetGroup()), maxGroupWidth))
		maxHost = maxInt(maxHost, minInt(getDisplayWidth(svr.GetHost()), maxHostWidth))
		maxUser = maxInt(maxUser, minInt(getDisplayWidth(svr.GetUser()), maxUserWidth))
	}

	// Calculate total width and print header
	totalWidth := 5 + maxName + maxGroup + maxHost + maxUser + 12 // ID + NAME + GROUP + HOST + USER + STATUS + spaces
	fmt.Println()
	printColoredLine("═", totalWidth, "cyan")
	fmt.Println()

	// Print table header
	printPadded("ID", 5, "")
	printPadded("NAME", maxName, "")
	printPadded("GROUP", maxGroup, "")
	printPadded("HOST", maxHost, "")
	printPadded("USER", maxUser, "")
	fmt.Println("STATUS")
	printColoredLine("─", totalWidth, "dim")
	fmt.Println()

	// Track groups for separators and stats
	var lastGroup string
	groupCounts := make(map[string]int)

	// Print rows
	for index, svr := range servers {
		group := svr.GetGroup()
		groupCounts[group]++

		// Add group separator when group changes
		if lastGroup != "" && lastGroup != group {
			printColoredLine("─", totalWidth, "dim")
			fmt.Println()
		}
		lastGroup = group

		// Print ID with cyan color
		printColoredPadded(fmt.Sprintf("[%d]", index), 5, colorCyan)

		// Print NAME with yellow color (truncate to maxNameWidth)
		printColoredPadded(truncate(svr.GetName(), maxNameWidth), maxName, colorYellow)

		// Print GROUP with blue color (truncate to maxGroupWidth)
		printColoredPadded(truncate(svr.GetGroup(), maxGroupWidth), maxGroup, colorBlue)

		// Print HOST (truncate to maxHostWidth)
		printPadded(truncate(svr.GetHost(), maxHostWidth), maxHost, "")

		// Print USER (red for root, truncate to maxUserWidth)
		userDisplay := svr.GetUser()
		if strings.EqualFold(userDisplay, "root") {
			printColoredPadded(truncate(userDisplay, maxUserWidth), maxUser, colorRed)
		} else {
			printPadded(truncate(userDisplay, maxUserWidth), maxUser, "")
		}

		// Print STATUS - check actual connection via ControlMaster if enabled
		if GetConf().GetSSHConfig().ShouldUseSystemSSH() {
			fmt.Println(GetConnectionStatus(svr.GetHost(), svr.GetPort(), svr.GetUser()))
		} else {
			fmt.Println(colorDim("○ Idle"))
		}
	}

	// Bottom separator
	printColoredLine("─", totalWidth, "dim")
	fmt.Println()

	// Print stats
	var stats []string
	stats = append(stats, fmt.Sprintf("Total: %d servers", len(servers)))
	for group, count := range groupCounts {
		stats = append(stats, fmt.Sprintf("%s: %d", group, count))
	}
	fmt.Println(colorDim(strings.Join(stats, " | ")))

	fmt.Println()
}

// getDisplayWidth returns the display width of a string
// Chinese characters take 2 display widths, ASCII take 1
func getDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r < 128 {
			width += 1
		} else {
			width += 2 // Chinese and other wide characters
		}
	}
	return width
}

// truncate truncates a string if it exceeds maxLen (by display width)
func truncate(s string, maxLen int) string {
	if getDisplayWidth(s) <= maxLen {
		return s
	}

	// Truncate by display width
	result := ""
	currentWidth := 0
	for _, r := range s {
		charWidth := 1
		if r >= 128 {
			charWidth = 2
		}
		if currentWidth+charWidth > maxLen-3 {
			break
		}
		result += string(r)
		currentWidth += charWidth
	}
	return result + "..."
}

// printPadded prints a string with padding (no color)
func printPadded(s string, width int, _ string) {
	displayWidth := getDisplayWidth(s)
	padding := width - displayWidth
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s ", s, strings.Repeat(" ", padding))
}

// printColoredPadded prints a colored string with proper padding
func printColoredPadded(s string, width int, colorFunc func(string) string) {
	colored := colorFunc(s)
	displayWidth := getDisplayWidth(s)
	padding := width - displayWidth
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s ", colored, strings.Repeat(" ", padding))
}

// maxInt returns the maximum of two integers
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper functions for colors
func colorRed(s string) string {
	return fmt.Sprintf("\033[1;31m%s\033[0m", s)
}

func colorYellow(s string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", s)
}

func colorBlue(s string) string {
	return fmt.Sprintf("\033[34m%s\033[0m", s)
}

func colorCyan(s string) string {
	return fmt.Sprintf("\033[36m%s\033[0m", s)
}

func colorDim(s string) string {
	return fmt.Sprintf("\033[2m%s\033[0m", s)
}

func colorGreen(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func printColoredLine(char string, width int, colorType string) {
	line := strings.Repeat(char, width)
	switch colorType {
	case "cyan":
		fmt.Print(colorCyan(line))
	case "dim":
		fmt.Print(colorDim(line))
	default:
		fmt.Print(line)
	}
}

func (c Config) Print() {
	c.PrintServerV2()
	fmt.Println()
}

func (c Config) Validate() error {
	if util.IsZero(c) {
		return fmt.Errorf("config is nil")
	}
	if len(c.GetServerList()) == 0 {
		return fmt.Errorf("server is empty")
	}
	for _, server := range c.ServerList {
		if err := server.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// readConfig reads the configuration from the config file
func readConfig() (*Config, error) {
	viper.SetConfigName("easy_config")
	viper.AddConfigPath(getConfigDir())
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	config := new(Config)
	err := viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	u, _ := user.Current()
	dir := fmt.Sprintf("%s/.easyssh", u.HomeDir)
	if !util.Exists(dir) {
		_ = util.Mkdir(dir)
	}
	return dir
}

// CheckSSHConnection checks if there's an active SSH connection via ControlMaster
// Returns true if connected, false otherwise
func CheckSSHConnection(host string, port int, user string) bool {
	// Use ssh -O check to see if there's an active connection
	// Address format: user@host (port via -p flag)
	args := []string{"-O", "check"}
	if port != 22 {
		args = append(args, "-p", strconv.Itoa(port))
	}
	args = append(args, fmt.Sprintf("%s@%s", user, host))

	cmd := exec.Command("ssh", args...)
	err := cmd.Run()
	return err == nil
}

// GetConnectionStatus returns a status string for display
func GetConnectionStatus(host string, port int, user string) string {
	if CheckSSHConnection(host, port, user) {
		return colorGreen("● Connected")
	}
	return colorDim("○ Idle")
}
