package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logging struct {
		StdOutPanics bool `yaml:"stdOutPanics"`
		HTTPVerbose  bool `yaml:"httpVerbose"`
	} `yaml:"logging"`
	Server struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Name string `yaml:"name"`
		URI  string `yaml:"uri"`
		Port int64  `yaml:"port"`
	} `yaml:"database"`
}

func RenderConfig() (*Config, error) {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
