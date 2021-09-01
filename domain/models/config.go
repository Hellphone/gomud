package models

type Config struct {
	Server struct {
		Port     string `yaml:"port"`
		Protocol string `yaml:"protocol"`
		Host     string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username string `yaml:"user"`
		Password string `yaml:"pass"`
	} `yaml:"database"`
}
