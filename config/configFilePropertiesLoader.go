package config

import (
	"fmt"
	"github.com/spf13/viper"

)
type ConfigFilePropertiesLoader struct {
	Properties Configurations
}	

func(c *ConfigFilePropertiesLoader) AppConnfig() ServerConfigurations {
	return c.Properties.Server
}

func(c *ConfigFilePropertiesLoader) DbConfig() DatabaseConfigurations {
	return c.Properties.Database
}

func NewConfigFilePropertiesLoader(filePath string) (PropertiesLoader, error) {
	viper.SetConfigFile(filePath)
	err1 := viper.ReadInConfig()
	if err1 != nil {
		return nil, fmt.Errorf("error while reading config file: %s, %w", filePath, err1)
		
	}

	var configuration Configurations

	err2 := viper.Unmarshal(&configuration)
	if err2 != nil {
		fmt.Printf("unable to decode into struct, %v", err2)
	}

	filePropertiesLoader := &ConfigFilePropertiesLoader{Properties: configuration}
	
	return filePropertiesLoader, nil
}