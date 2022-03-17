package transact

import (
	"bufio"
	"fmt"
	"go_storage/config"
	"go_storage/core"
	"os"

)


type FileTransactionLogger struct {
	events 		 chan<- core.Event	
	errors       <-chan error
	lastSequence uint64
	file         *os.File
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{ EventType: core.EventPut, Key: key, Value: value}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{ EventType: core.EventDelete, Key: key}
}

func (l *FileTransactionLogger) Err() <- chan error {
	return l.errors
}

func (l *FileTransactionLogger) LastSequence() uint64 {
	return l.lastSequence
}

func (l *FileTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan core.Event)
	outError := make(chan error, 1)

	go func() {
		var e core.Event 
		defer close(outEvent)
		defer close(outError)
		
		for scanner.Scan() {
			line := scanner.Text() 

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s", &e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}
			
			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}
			l.lastSequence = e.Sequence
			outEvent <- e
		}
		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
			return
		}
	}()
	return outEvent, outError
}

func (l *FileTransactionLogger) Run() {
	events := make(chan core.Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		for e := range events {
			l.lastSequence ++
			_, err := fmt.Fprintf(l.file, "%d\t%d\t%s\t%s\n", l.lastSequence, e.EventType, e.Key, e.Value)
			
			if err != nil {	
				errors <- err
				return
		}
	}
	}()
}

func NewFileTransactionLogger(config config.ServerConfigurations) (core.TransactionLogger, error) {
	file, err := os.OpenFile(config.LogFile, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0755)
	
	if err != nil {
		fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file}, nil 

}
