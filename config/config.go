package config

import (
	"github.com/spf13/viper"
	"log"
)

type OperatorConfig struct {
	Version                   string `mapstructure:"VERSION"`
	DefaultOpenNMSValuesFile  string `mapstructure:"DEFAULT_OPENNMS_VALUES_FILE"`
	DefaultOperatorValuesFile string `mapstructure:"DEFAULT_OPERATOR_VALUES_FILE"`
	DevMode                   bool   `mapstructure:"DEV_MODE"`

	//Image update
	ImageUpdateFreq int `mapstructure:"IMAGE_UPDATE_FREQUENCY"`

	//Service and job images
	ServiceImageGrafana string `mapstructure:"SERVICE_IMAGE_GRAFANA"`

	//instance node restrictions
	NodeRestrictionKey   string `mapstructure:"NODE_RESTRICTION_KEY"`
	NodeRestrictionValue string `mapstructure:"NODE_RESTRICTION_VALUE"`
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
