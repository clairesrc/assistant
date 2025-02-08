package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	ollamaClient, err := newOllamaClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create Ollama API client: %w", err))
	}

	// set up web server
	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		err := getUpdates(w, r, ollamaClient)
		writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
	})

	http.ListenAndServe(":8080", nil)
}

func getUpdates(w http.ResponseWriter, r *http.Request, o *ollamaClient) error {
	// validate POST request
	if r.Method != http.MethodPost {
		return fmt.Errorf("method not allowed")
	}

	// get prompts from request body
	var prompts []string
	err := json.NewDecoder(r.Body).Decode(&prompts)
	if err != nil {
		return fmt.Errorf("cannot decode request body: %w", err)
	}

	// concurrently generate updates for each prompt
	updates := make([]string, len(prompts))
	var wg sync.WaitGroup
	for i, prompt := range prompts {
		wg.Add(1)
		go func(i int, prompt string) {
			defer wg.Done()
			updates[i], err = o.generate(prompt)
			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
			}
		}(i, prompt)
	}
	wg.Wait()

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
