package frontend

import "go_storage/core"

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error 
}

