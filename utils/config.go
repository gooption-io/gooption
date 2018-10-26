package utils

import (
	"github.com/spf13/viper"
)

/*
This package uses viper to load config from files

Examples:

**Simple config**

config.yaml:

tcp: :50051
http: :8081
promhttp: :8080

golang:

config := NewServiceConfig("")

-----

**Env specific config**

config.yaml:

dev:
  tcp: :1003
  http: :1002
  promhttp: :1001

prod:
  tcp: :50051
  http: :8081
  promhttp: :8080

golang:

config := NewServiceConfig("dev")

*/

// ServiceConfig contains ports the service will listen on
type ServiceConfig struct {
	TCP      string `mapstructure:"tcp"`
	HTTP     string `mapstructure:"http"`
	PromHTTP string `mapstructure:"promhttp"`
}

// InitViperConfig reads viper config
// viper tries all paths
func InitViperConfig(name string, paths ...string) {
	viper.SetConfigName(name)
	for _, path := range paths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// NewServiceConfig creates new config from yaml
// env must match the config file
func NewServiceConfig(root string) ServiceConfig {
	var config ServiceConfig
	if root == "" {
		err := viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}
		return config
	}

	err := viper.UnmarshalKey(root, &config)
	if err != nil {
		panic(err)
	}
	return config
}
