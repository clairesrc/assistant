package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type calendarClient interface {
	getEvents() ([]calendarEvent, error)
}

type calendar struct {
	apiKey string
	baseURL string
}

type calendarEvent struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Start string `json:"start"`
	End string `json:"end"`
}

func newCalendarClient() (calendarClient, error) {
	if os.Getenv("CALENDAR_API_KEY") == "" {
		return nil, fmt.Errorf("CALENDAR_API_KEY is not set")
	}

	if os.Getenv("CALENDAR_BASE_URL") == "" {
		return nil, fmt.Errorf("CALENDAR_BASE_URL is not set")
	}

	return &calendar{
		apiKey: os.Getenv("CALENDAR_API_KEY"),
		baseURL: os.Getenv("CALENDAR_BASE_URL"),
	}, nil
}

func (c *calendar) getEvents() ([]calendarEvent, error) {
	url := fmt.Sprintf("%s/events?apiKey=%s", c.baseURL, c.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get events: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read events: %w", err)
	}

	var events []calendarEvent
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal events: %w", err)
	}

	return events, nil
}
