package config

import "time"

type SSHConfig struct {
	Key               string `json:"key" yaml:"key"`
	KeepAlive         bool   `json:"keep_alive" yaml:"keep_alive"`
	KeepAliveInterval string `json:"keep_alive_interval" yaml:"keep_alive_interval"` // e.g., "60s"
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
