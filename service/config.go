package service

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config represents service configuration.
type Config struct {
	DataDir  string `yaml:"dataDir"`
	Network  string `yaml:"network"`
	Address  string `yaml:"address"`
	DirDepth int    `yaml:"dirDepth"`
}

// LoadConfig loads the given configuration file.
func LoadConfig(filename string) (*Config, error) {
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var config *Config

	return config, yaml.Unmarshal(b, &config)
}
