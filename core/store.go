package core

import (
	"errors"
	"log"
	"sync"
)

// ErrorNoSuchKey thrown when key not found
var ErrorNoSuchKey = errors.New("no such key")

var store = struct{
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

// Put key in store
func Put(key, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

// Get value form store
func Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

// Delete key in store
func Delete(key string) error {
	store.Lock()
	_, ok := store.m[key]
	store.Unlock()
	if !ok {
		return ErrorNoSuchKey
	} 
	delete(store.m, key)
	return nil
}

func (store *KeyValueStore) Restore() error {
	var err error

	events, errors := store.transact.ReadEvents()
	count, ok, e := 0, true, Event{}
	
		
	for ok && err == nil {

		select {
			case err, ok = <-errors:
			log.Printf("1")
		
			case e, ok = <-events:
				log.Printf("1")
				switch e.EventType {
					
					case EventDelete: // Got a DELETE event!
						log.Printf("2")
						err = store.Delete(e.Key)
						count++
					case EventPut: // Got a PUT event!
						log.Printf("3")
						err = store.Put(e.Key, e.Value)
						count++
				}
		}
		log.Printf("2")
	}

	log.Printf("%d events replayed\n", count)

	store.transact.Run()

	go func() {
		for err := range store.transact.Err() {
			log.Print(err)
		}
	}()

	return err
}