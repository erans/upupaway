package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// PathConfig provides the configuration for a single path in the system
type PathConfig struct {
	Path                      string `yaml:"path"`
	BucketName                string `yaml:"bucketName"`
	AccessControlAllowOrigin  string `yaml:"accessControlAllowOrigin"`
	AccessControlAllowMethods string `yaml:"accessControlAllowMethods"`
	AccessControlMaxAge       string `yaml:"accessControlMaxAge"`
}

// Config provides a configuration struct for the server
type Config struct {
	DebugLevel string                   `yaml:"debugLevel"`
	Buckets    []map[string]interface{} `yaml:"buckets"`
	Paths      []PathConfig             `yaml:"paths"`
	Storage    struct {
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

var cfg *Config

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

// SetGlobalConfig sets the global config
func SetGlobalConfig(c *Config) {
	cfg = c
}

// GetGlobalConfig gets the global config
func GetGlobalConfig() *Config {
	return cfg
}

var pathConfigIndex = map[string]PathConfig{}

// GetPathConfigByPath returns the path config by the given path
func GetPathConfigByPath(path string) *PathConfig {
	if len(pathConfigIndex) == 0 {
		for _, p := range GetGlobalConfig().Paths {
			pathConfigIndex[p.Path] = p
		}
	}

	if c, ok := pathConfigIndex[path]; ok {
		return &c
	}

	return nil
}
