package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":8080"

var store = make(map[string]string)
var ErrorNoSuchKey = errors.New("no such key")

func main() {
	log.Printf("service is running on %s port", port)
	router := mux.NewRouter()
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	log.Fatal(http.ListenAndServe(port, router))
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
	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Put(key, value string) error {
	store[key] = value

	return nil
}

func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

func Delete(key string) error {
	delete(store, key)

	return nil
}
