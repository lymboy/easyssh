package config

import (
	"easyssh/util"
	"fmt"
	"os"
	"os/user"
	"sort"
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
	// Print header
	fmt.Println()
	printColoredLine("═", 80, "cyan")
	fmt.Println()

	// Print table header
	headerFmt := "%-5s %-18s %-12s %-18s %-12s %s"
	fmt.Printf(headerFmt, "ID", "NAME", "GROUP", "HOST", "USER", "STATUS")
	fmt.Println()
	printColoredLine("─", 80, "dim")
	fmt.Println()

	// Track groups for separators and stats
	var lastGroup string
	groupCounts := make(map[string]int)

	// Print rows
	for index, svr := range c.GetServerList() {
		group := svr.GetGroup()
		groupCounts[group]++

		// Add group separator when group changes
		if lastGroup != "" && lastGroup != group {
			printColoredLine("─", 80, "dim")
			fmt.Println()
		}
		lastGroup = group

		// Format user with color (red for root)
		userDisplay := svr.GetUser()
		if strings.EqualFold(userDisplay, "root") {
			userDisplay = colorRed(userDisplay)
		}

		// Format status
		statusDisplay := colorDim("○ Idle")

		// Print row
		fmt.Printf("%-5s %-18s %-12s %-18s %-12s %s",
			colorCyan(fmt.Sprintf("[%d]", index)),
			colorYellow(svr.GetName()),
			colorBlue(group),
			svr.GetHost(),
			userDisplay,
			statusDisplay,
		)
		fmt.Println()
	}

	// Bottom separator
	printColoredLine("─", 80, "dim")
	fmt.Println()

	// Print stats
	var stats []string
	stats = append(stats, fmt.Sprintf("Total: %d servers", len(c.GetServerList())))
	for group, count := range groupCounts {
		stats = append(stats, fmt.Sprintf("%s: %d", group, count))
	}
	fmt.Println(colorDim(strings.Join(stats, " | ")))

	fmt.Println()
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
