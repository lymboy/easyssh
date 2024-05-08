package config

import (
	"easyssh/util"
	"fmt"
	"os"
	"os/user"
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
	SSHConfig SSHConfig `json:"ssh" yaml:"ssh" mapstructure:"ssh"`
	Server    []Server  `json:"server" yaml:"server"`
}

func (c Config) GetSSHConfig() *SSHConfig {
	if util.IsZero(c.SSHConfig) {
		return new(SSHConfig)
	}
	return &c.SSHConfig
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
	for _, server := range c.Server {
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
	if index < len(c.Server) {
		return &c.Server[index]
	}
	return nil
}

func (c Config) GetServerForMap() map[string]*Server {
	serverMap := make(map[string]*Server)
	for _, server := range c.Server {
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
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>  Server  <<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println()
	for i, server := range c.Server {
		fmt.Printf("\t[%d]\t%s\t\t%s\n", i, server.GetName(), server)
	}
}

func (c Config) Print() {
	fmt.Println("***********************  EASY SSH  ************************")
	c.PrintServer()
	fmt.Println()
	fmt.Println("***********************************************************")
}

func (c Config) Validate() error {
	if util.IsZero(c) {
		return fmt.Errorf("config is nil")
	}
	if len(c.Server) == 0 {
		return fmt.Errorf("server is empty")
	}
	for _, server := range c.Server {
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
		util.Mkdir(dir)
	}
	return dir
}
