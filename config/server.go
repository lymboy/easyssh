package config

import (
	"easyssh/util"
	"fmt"
	"os/user"
)

type Server struct {
	Name     string `json:"name" yaml:"name"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Parent   string `json:"parent" yaml:"parent"`
	Desc     string `json:"desc" yaml:"desc"`
}

func (s Server) GetName() string {
	if len(s.Name) == 0 {
		return "Default"
	}
	return s.Name
}

func (s Server) GetHost() string {
	return s.Host
}

func (s Server) GetPort() int {
	if s.Port <= 0 {
		return 22
	}
	return s.Port
}

func (s Server) GetUser() string {
	if len(s.User) == 0 {
		u, _ := user.Current()
		return u.Username
	}
	return s.User
}

func (s Server) GetPassword() string {
	return s.Password
}

func (s Server) GetParent() string {
	return s.Parent
}

func (s Server) GetDesc() string {
	return s.Desc
}

func (s Server) String() string {
	return fmt.Sprintf("%s@%s", s.GetUser(), s.Host)
}

func (s Server) GenSSHCommand() string {
	return fmt.Sprintf("ssh %s", s.String())
}

func (s Server) GenSpawnCommand() string {
	return fmt.Sprintf("spawn %s", s.GenSSHCommand())
}

func (s Server) GenSCPCommand(src, dest string) string {
	return fmt.Sprintf("scp %s %s:%s", src, s.String(), dest)
}

func (s Server) Validate() error {
	if util.IsZero(s) {
		return fmt.Errorf("server is empty")
	}
	if len(s.Host) == 0 {
		return fmt.Errorf("host is not set")
	}
	return nil
}
