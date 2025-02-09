package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type prompt struct {
	key string
	prompt string
	generateImage bool
}

// PromptResult is the response value for a given prompt
type PromptResult struct {
	Key string `json:"key"`
	Prompt prompt `json:"prompt"`
	ImageURL string `json:"image_url,omitempty"`
}

func main() {
	openWebUIClient, err := newOpenWebUIClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create openWebUI API client: %w", err))
	}

	automaticSDClient, err := newAutomaticSDClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create automaticSD API client: %w", err))
	}

	// set up web server
	http.HandleFunc("/updates", func(w http.ResponseWriter, r *http.Request) {
		err := getUpdates(w, r, openWebUIClient, automaticSDClient)
		writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
	})

	http.ListenAndServe(":8080", nil)
}

func getUpdates(w http.ResponseWriter, r *http.Request, o *openWebUIClient, a *automaticSDClient) error {
	// get source data: weather
	weather, err := getWeather()
	if err != nil {
		return fmt.Errorf("cannot get weather: %w", err)
	}

	// get source data: news
	news, err := getNews()
	if err != nil {
		return fmt.Errorf("cannot get news: %w", err)
	}

	// get source data: calendar events
	calendarEvents, err := getCalendarEvents()
	if err != nil {
		return fmt.Errorf("cannot get calendar events: %w", err)
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
