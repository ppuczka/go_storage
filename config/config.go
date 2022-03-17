package config


type Configurations struct {
	Server       ServerConfigurations
	Database     DatabaseConfigurations

}

type DatabaseConfigurations struct {
	DbName     string
	DbHost     string
	DbUser     string
	DbPassword string

}

type ServerConfigurations struct {
	AppPort      string
	AppHost      string
	LogType      string
	LogFile		 string
	FrontEndType string
	TLSCert      string
	PrivateKey   string
	
}