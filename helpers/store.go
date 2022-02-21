package helpers

import (
	"errors"
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
