package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
}

// ComicConfig struct contains the comic configuration
type ComicConfig struct {
	DefaultOutputFormat string `mapstructure:"default_output_format"`
}

// LoadConfig read the `config` file and unmarshal to struct
func (c *ComicConfig) LoadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(&c)
}
