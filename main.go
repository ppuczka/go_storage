package main

import (
	"errors"
	"fmt"
	"go_storage/helpers"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":8080"
var logger helpers.TransactionLogger

func main() {
	properties, err := helpers.NewEnvironmentVarsPropertiesLoader()
	if err != nil {
		log.Fatal("error while loading application properties from env variables: ", err)
	}

	initializeTransactionLog(properties)
	
	router := mux.NewRouter()
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueReadHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	
	_, listenPort := properties.AppConnfig()
	log.Fatal(http.ListenAndServe(listenPort, router))
	log.Printf("service is running on %s port", port)

}

func initializeTransactionLog(properties helpers.PropertiesLoader) error {
	log.Printf("Initializing DB connection")
	var err error 
	logger, err = helpers.NewPostgresTransactionLogger(properties.DbConfig())
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	log.Printf("event logger initialized")
	
	events, errors := logger.ReadEvents()
	e, ok := helpers.Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <- errors:
		case e, ok = <- events:
			switch e.EventType {
			case helpers.EventDelete:
				err = helpers.Delete(e.Key)
			case helpers.EventPut: 
			err = helpers.Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()
	return err
}

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := io.ReadAll(r.Body)
	r.Body.Close()
	log.Printf("recived PUT request with key: %s and value: %s", key, value)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = helpers.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(value))
	log.Printf("---- Created ----")
	w.WriteHeader(http.StatusCreated)
}

func keyValueReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := helpers.Get(key)
	log.Printf("recived GET request with key: %s", key)
	
	if errors.Is(err, helpers.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value))
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	err := helpers.Delete(key)

	log.Printf("recived DELETE request with key: %s", key)

	if errors.Is(err, helpers.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	logger.WriteDelete(key)
	log.Printf("---- DELETED ----")

	w.WriteHeader(http.StatusCreated)
}

