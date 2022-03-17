package transact

import (
	"fmt"
	"go_storage/config"
	"go_storage/core"
)

func NewTransactionLogger(properties config.PropertiesLoader) (core.TransactionLogger, error) {
	switch properties.AppConnfig().LogType {
	case "file":
		return NewFileTransactionLogger(properties.AppConnfig())
	
	case "db":
		return NewPostgresTransactionLogger(properties.DbConfig())
	
	default:
		return nil, fmt.Errorf("no such transaction logger")
	}

}
