package helpers

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const tableName = "storage"

type PostgresTransactionLogger struct {
	events      chan <- Event
	errors      <- chan error
	db          *sql.DB
}


func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- Event{ EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WrieDelete(key string) {
	l.events <- Event{ EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) Err() <- chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", tableName))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != tableName {
		rows.Scan(&result)
	}

	return result == tableName, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error
	
	createQuery := `CREATE TABLE transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err 
	}

	return nil
}

func NewPostgresTransactionLogger(config PostrgesDBParams) (TransactionLogger, error) {

	connString := fmt.Sprintf("host=%s, dbname=%s, user=%s, password=%s", config.host, config.dbName, config.user, config.password)
	db, err := sql.Open("postgress", connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to open db: %w", err)
	}
	
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to establish db connection %w", err)
	}

	logger := &PostgresTransactionLogger{db: db}

	exists, err := logger.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}

	if !exists {
		if err = logger.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		
		}
	}

	return logger, nil

}
