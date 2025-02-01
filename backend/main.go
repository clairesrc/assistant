package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ollamaClient, err := newOllamaClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create Ollama API client: %w", err))
	}

	// set up web server
	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		err := getUpdates(w, r, ollamaClient)
		log.Default().Println(fmt.Errorf("cannot get updates: %w", err))
	})

	http.ListenAndServe(":8080", nil)
}

func getUpdates(w http.ResponseWriter, _ *http.Request, o *ollamaClient) error {
	updates, err := o.generate("What's the weather like today?")
	if err != nil {
		writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
	}

	// return updates as json
	updatesJson, err := json.Marshal(updates)
	if err != nil {
		writeHttpError(w, http.StatusInternalServerError, "cannot marshal updates to json", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(updatesJson)

	return nil
}

func writeHttpError(w http.ResponseWriter, status int, message string, err error) {
	fmt.Println(fmt.Errorf("%s: %w", message, err))
	w.WriteHeader(status)
	w.Write([]byte(message))
}
