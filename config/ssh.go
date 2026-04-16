package config

import "time"

type SSHConfig struct {
	Key               string `json:"key" yaml:"key" mapstructure:"key"`
	KeepAlive         bool   `json:"keep_alive" yaml:"keep_alive" mapstructure:"keep_alive"`
	KeepAliveInterval string `json:"keep_alive_interval" yaml:"keep_alive_interval" mapstructure:"keep_alive_interval"` // e.g., "60s"
	UseSystemSSH      bool   `json:"use_system_ssh" yaml:"use_system_ssh" mapstructure:"use_system_ssh"`                 // Use system ssh command (enables ControlMaster)
}

func (s SSHConfig) GetKey() string {
	if len(s.Key) == 0 {
		return "id_rsa"
	}
	return s.Key
}

func (s SSHConfig) GetKeepAliveInterval() time.Duration {
	if s.KeepAliveInterval == "" {
		return 60 * time.Second // default
	}
	d, err := time.ParseDuration(s.KeepAliveInterval)
	if err != nil {
		return 60 * time.Second
	}
	return d
}

// ShouldUseSystemSSH returns true if system SSH should be used
// When enabled, delegates to system ssh command which supports ControlMaster
func (s SSHConfig) ShouldUseSystemSSH() bool {
	return s.UseSystemSSH // User can explicitly enable, or we could auto-detect
}
