package config

type PropertiesLoader interface {

	DbConfig() DatabaseConfigurations
	AppConnfig() ServerConfigurations

}