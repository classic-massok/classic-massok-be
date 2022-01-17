package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int64  `yaml:"port"`
	} `yaml:"server"`
	Logging struct {
		StdOutPanics bool `yaml:"stdOutPanics"`
		HTTPVerbose  bool `yaml:"httpVerbose"`
	} `yaml:"logging"`
	Database struct {
		Name string `yaml:"name"`
		URI  string `yaml:"uri"`
		Port int64  `yaml:"port"`
	} `yaml:"database"`
	Tokens struct {
		AccessTokenPrvKeyPath  string `yaml:"accessTokenPrvKeyPath"`
		AccessTokenPubKeyPath  string `yaml:"accessTokenPubKeyPath"`
		RefreshTokenPrvKeyPath string `yaml:"refreshTokenPrvKeyPath"`
		RefreshTokenPubKeyPath string `yaml:"refreshTokenPubKeyPath"`
	} `yaml:"tokens"`
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
