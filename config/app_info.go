package config

type AppInfo struct {
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Author      string `json:"author" yaml:"author"`
	Description string `json:"description" yaml:"description"`
}
