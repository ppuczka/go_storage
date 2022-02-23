package helpers

import (
	"fmt"
	"os"
)

type EnvironmentVarsPropertiesLoader struct {
	Properties AppProperties
}

func(e *EnvironmentVarsPropertiesLoader) DbConfig() PostrgesDBParams {
	
	dbParams := PostrgesDBParams{DbName: e.Properties.DatabaseName, Host: e.Properties.DatabaseHost, 
					User: e.Properties.DatabaseUser, Password: e.Properties.DatabasePass}
	
	return dbParams
}

func (e *EnvironmentVarsPropertiesLoader) AppConnfig() (host string, port string) {
	return e.Properties.AppHost, e.Properties.AppPort
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

	loadedProperties := AppProperties{AppPort: port, AppHost: host,  DatabaseHost: dbHost, DatabaseUser: dbUser, DatabasePass: dbPass, DatabaseName: dbName}
	envPropertiesLoader := &EnvironmentVarsPropertiesLoader{Properties: loadedProperties}
	
	return envPropertiesLoader, nil
}