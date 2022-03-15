package transact

import (
	"database/sql"
	"fmt"
	"go_storage/config"
	"go_storage/core"
	"log"

	_ "github.com/lib/pq"
)

const tableName = "transactions"

type PostgresTransactionLogger struct {
	events      chan <- core.Event
	errors      <- chan error
	db          *sql.DB
}


func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{ EventType: core.EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{ EventType: core.EventDelete, Key: key}
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

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions (event_type, key, value) VALUES ($1, $2, $3)`

		for e := range events {
			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
		
			if err != nil {
				errors <- err
			}	
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()
		e := core.Event{}

		for rows.Next() {

			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
				return
			}
			outEvent <- e
		}
		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()
	return outEvent, outError
}

func NewPostgresTransactionLogger(config config.DatabaseConfigurations) (core.TransactionLogger, error) {

	connString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable" , config.DbHost, config.DbName, config.DbUser, config.DbPassword)
	
	log.Printf("---- opening connection with db %s %s ----", config.DbHost, config.DbName)
	
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}
	log.Printf("---- connection with db %s %s opened ----", config.DbHost, config.DbName)
	
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to establish db connection %w", err)
	}
	log.Printf("---- PING sucessfull ----")
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
