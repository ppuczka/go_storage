package helpers

type PropertiesLoader interface {

	DbConfig() PostrgesDBParams
	AppConnfig() (string, string)

}