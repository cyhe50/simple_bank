package util

import (
	"github.com/spf13/viper"
)

type EnvConfig struct {
	DriverName     string `mapstructure:"DRIVER_NAME"`
	DataSourceName string `mapstructure:"DATA_SOURCE_NAME"`
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config EnvConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
