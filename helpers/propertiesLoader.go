package helpers

import(
	"go_storage/config"
)

type PropertiesLoader interface {

	DbConfig() config.DatabaseConfigurations
	AppConnfig() config.ServerConfigurations

}