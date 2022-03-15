package core

import(

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
	delete(store.m, key)
	store.transact.WriteDelete(key)
	return nil
}

func (store *KeyValueStore) Put(key string, value string) error {
	store.m[key] = value
	store.transact.WritePut(key, value)
	return nil
}