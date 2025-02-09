package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// remote source data clients
	weatherClient, err := newWeatherClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create weather client: %w", err))
	}

	newsClient, err := newNewsClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create news client: %w", err))
	}

	calendarClient, err := newCalendarClient()
	if err != nil {
		log.Fatal(fmt.Errorf("can't create calendar client: %w", err))
	}

	// local ai clients
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
		err := getUpdates(w, r, openWebUIClient, automaticSDClient, weatherClient, newsClient, calendarClient)
		writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
	})

	http.ListenAndServe(":8080", nil)
}

func writeHttpError(w http.ResponseWriter, status int, message string, err error) {
	fmt.Println(fmt.Errorf("%s: %w", message, err))
	w.WriteHeader(status)
	w.Write([]byte(message))
}
