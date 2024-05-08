package config

type SSHConfig struct {
	Key string `json:"key" yaml:"key"`
	KeepAlive bool `json:"keep_alive" yaml:"keep_alive"`
}

func (s SSHConfig) GetKey() string {
	if len(s.Key) == 0 {
		return "id_rsa"
	}
	return s.Key
}
