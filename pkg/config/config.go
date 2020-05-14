package config

import (
	"bytes"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PublicKeysDirectory string `yaml:"public_keys_directory" mapstructure:"public_keys_directory"`
	HomeContent         string `yaml:"home_content" mapstructure:"home_content"`
}

var (
	defaultConfig = Config{
		PublicKeysDirectory: "",
		HomeContent:         "",
	}
)

func Read() *Config {
	config := &Config{}
	viper.SetConfigType("yaml")

	defaultConfigYaml, _ := yaml.Marshal(defaultConfig)
	viper.ReadConfig(bytes.NewBuffer(defaultConfigYaml))

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	viper.Unmarshal(config)
	return config
}
