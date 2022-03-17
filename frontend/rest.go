package frontend

import (
	"errors"
	"go_storage/config"
	"go_storage/core"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type restFrontEnd struct {
	store      *core.KeyValueStore
	properties *config.ServerConfigurations 
}

func (f *restFrontEnd) Start(store *core.KeyValueStore) error {
	store = f.store
	router := mux.NewRouter()
	
	router.HandleFunc("/v1/{key}", f.keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", f.keyValueReadHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", f.keyValueDeleteHandler).Methods("DELETE")
	
	listenPort := f.properties.AppPort
	
	log.Printf("service is running on %s port", listenPort)

	return http.ListenAndServeTLS(listenPort, f.properties.TLSCert, f.properties.PrivateKey, router)
}


func (f *restFrontEnd) keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	
	vars := mux.Vars(r)
	key := vars["key"]

	err := f.store.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (f *restFrontEnd) keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["value"]

	err := f.store.Put(key, val)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (f *restFrontEnd) keyValueReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := f.store.Get(key)
	if errors.Is(err, core.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))

	log.Printf("GET key=%s\n", key)
}