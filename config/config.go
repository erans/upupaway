package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config provides a configuration struct for the server
type Config struct {
	DebugLevel string                   `yaml:"debugLevel"`
	Buckets    []map[string]interface{} `yaml:"buckets"`
	Paths      []struct {
		Path       string `yaml:"path"`
		BucketName string `yaml:"bucketName"`
	} `yaml:"paths"`
	Storage struct {
		ActiveStorage string `yaml:"activeStorage"`
		InMemory      struct {
			Size       int `yaml:"size"`
			Expiration int `yaml:"expiration"`
		} `yaml:"inmemory"`
		Redis struct {
			Host       string `yaml:"host"`
			Port       int    `yaml:"port"`
			DB         int    `yaml:"db"`
			Expiration int    `yaml:"expiration"`
		} `yaml:"redis"`
		MemCache struct {
			Host       string `yaml:"host"`
			Port       int    `yaml:"port"`
			Expiration int32  `yaml:"expiration"`
		} `yaml:"memcache"`
	} `yaml:"storage"`
}

// LoadConfig loads the config file
func LoadConfig(configFile string) (*Config, error) {
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configFile); err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		return nil, err
	}

	return &c, nil
}
