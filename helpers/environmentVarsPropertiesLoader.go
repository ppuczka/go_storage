package helpers

import (
	"fmt"
	"os"
	"go_storage/config"
)

type EnvironmentVarsPropertiesLoader struct {
	Properties config.Configurations
}

func(e *EnvironmentVarsPropertiesLoader) DbConfig() config.DatabaseConfigurations {
	return e.Properties.Database
}

func (e *EnvironmentVarsPropertiesLoader) AppConnfig() config.ServerConfigurations {
	return e.Properties.Server
}


func NewEnvironmentVarsPropertiesLoader() (PropertiesLoader, error) {
	fmt.Printf("Loading properties from env variables\n")
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("PORT env not present")	
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		return nil, fmt.Errorf("HOST env not present")	
	}
	
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, fmt.Errorf("DB_HOST env not present")	
	}

	dbUser, ok := os.LookupEnv("DB_USER")
	if !ok {
		return nil, fmt.Errorf("DB_USER env not present")	
	}

	dbPass, ok := os.LookupEnv("DB_PASS")
	if !ok {
		return nil, fmt.Errorf("DB_PASS env not present")	
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, fmt.Errorf("DB_NAME env not present")	
	}

	serverConfig := config.ServerConfigurations{AppPort: port, AppHost: host}
	dbConfig := config.DatabaseConfigurations{DbHost: dbHost, DbUser: dbUser, DbPassword: dbPass, DbName: dbName}
	loadedConfig := config.Configurations{serverConfig, dbConfig}

	envPropertiesLoader := &EnvironmentVarsPropertiesLoader{Properties: loadedConfig}
	return envPropertiesLoader, nil
}