package config

import (
	"github.com/spf13/viper"
	"log"
)

type OperatorConfig struct {
	Version                  string `mapstructure:"VERSION"`
	DefaultOpenNMSValuesFile string `mapstructure:"DEFAULT_OPENNMS_VALUES_FILE"`

	//Service and job images
	ServiceImageAuth    string `mapstructure:"SERVICE_IMAGE_AUTH"`
	ServiceImageGrafana string `mapstructure:"SERVICE_IMAGE_GRAFANA"`
	ServiceImageInit    string `mapstructure:"SERVICE_IMAGE_INIT"`
}

func LoadConfig() OperatorConfig {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	var config OperatorConfig

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err.Error())
	}
	return config
}
