package config

import (
	"easyssh/util"
	"fmt"
	"github.com/olekukonko/tablewriter"
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
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>  ServerList  <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println()
	// 创建表格
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Group", "Name", "Host", "Port", "User", "Password", "Parent", "Desc"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
	)
	table.SetCaption(true, "Welcome visit: https://lymboy.com")
	table.SetAutoMergeCells(true)
	table.SetAutoMergeCellsByColumnIndex([]int{1})
	//table.SetRowLine(true)
	// 添加数据到表格
	for index, svr := range c.GetServerList() {
		// 如果是 root 用户，加红显示
		if strings.EqualFold(svr.GetUser(), "root") {
			cols := []string{cast.ToString(index), svr.GetGroup(), svr.GetName(), svr.GetHost(), fmt.Sprintf("%d", svr.GetPort()), svr.GetUser(), svr.GetPassword(), svr.GetParent(), svr.GetDesc()}
			table.Rich(cols, []tablewriter.Colors{tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}})
			continue
		}
		table.Append([]string{cast.ToString(index), svr.GetGroup(), svr.GetName(), svr.GetHost(), fmt.Sprintf("%d", svr.GetPort()), svr.GetUser(), svr.GetPassword(), svr.GetParent(), svr.GetDesc()})
	}
	// 渲染表格
	table.Render()
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
