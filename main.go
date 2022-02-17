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

// Error when key not found 
var ErrorNoSuchKey = errors.New("no such key")

func main() {
	log.Printf("service is running on %s port", port)
	
	router := mux.NewRouter()
	router.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	router.HandleFunc("/v1/{key}", keyValueReadHandler).Methods("GET")
	router.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	
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
	log.Printf("---- Created ----")
	w.WriteHeader(http.StatusCreated)
}

func keyValueReadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value, err := Get(key)
	log.Printf("recived GET request with key: %s", key)
	
	if errors.Is(err, ErrorNoSuchKey) {
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
	err := Delete(key)

	log.Printf("recived DELETE request with key: %s", key)

	if errors.Is(err, ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	log.Printf("---- DELETED ----")

	w.WriteHeader(http.StatusCreated)
}

// Put key
func Put(key, value string) error {
	store[key] = value

	return nil
}

// Return value
func Get(key string) (string, error) {
	value, ok := store[key]
	if !ok {
		return "", ErrorNoSuchKey
	}
	return value, nil
}

// Delete key 
func Delete(key string) error {
	_, ok := store[key]
	if !ok {
		return ErrorNoSuchKey
	} 
	delete(store, key)
	return nil
}
