package core

import (
	"log"
	"sync"
)
type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut 
)

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)	
	Err() <-chan error

	ReadEvents() (<-chan Event, <-chan error)

	Run()
}

type KeyValueStore struct {
	sync.RWMutex
	m        map[string]string
	transact TransactionLogger 
}

func NewKeyValueStore(tl TransactionLogger) *KeyValueStore {
	return &KeyValueStore{
		m:        make(map[string]string),
		transact: tl,
	}
}


func (store *KeyValueStore) Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	log.Printf("11111")
	store.transact.WriteDelete(key)
	
	return nil
}

func (store *KeyValueStore) Get(key string) (string, error) {
	store.RLock()
	value, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return value, nil
}

func (store *KeyValueStore) Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	log.Printf("12345")
	store.transact.WritePut(key, value)

	return nil
}