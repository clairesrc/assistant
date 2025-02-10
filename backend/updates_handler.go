package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type prompt struct {
	key           string
	prompt        string
	generateImage bool
}

// PromptResult is the response value for a given prompt
type PromptResult struct {
	Key      string `json:"key"`
	Response string `json:"response"`
	ImageURL string `json:"image_url,omitempty"`
}

func getUpdates(w http.ResponseWriter, _ *http.Request, o openWebUIClient, a automaticSDClient, weather weatherClient, news newsClient, calendar calendarClient) error {
	// get source data: weather
	weatherResult, err := weather.get()
	if err != nil {
		return fmt.Errorf("cannot get weather: %w", err)
	}

	// get source data: news
	newsResults, err := news.get()
	if err != nil {
		return fmt.Errorf("cannot get news: %w", err)
	}

	// get source data: calendar events
	calendarEvents, err := calendar.getEvents()
	if err != nil {
		return fmt.Errorf("cannot get calendar events: %w", err)
	}

	prompts := []prompt{
		{
			key:           "weather",
			prompt:        fmt.Sprintf("You are a weather assistant. The current temperature is %f°C and the weather is %s. Write a very short comment on the weather.", weatherResult.Temp, weatherResult.Weather),
			generateImage: false,
		},
		{
			key:           "news",
			prompt:        fmt.Sprintf("You are a news assistant. The latest news are below: %v.\n \n Write a very short comment on the news.", newsResults),
			generateImage: true,
		},
		{
			key:           "calendar1",
			prompt:        fmt.Sprintf("You are a calendar assistant. The calendar event is below: %v.\n \n Write a very short comment on the calendar event.", calendarEvents[0]),
			generateImage: false,
		},
		{
			key:           "calendar2",
			prompt:        fmt.Sprintf("You are a calendar assistant. The calendar event is below: %v.\n \n Write a very short comment on the calendar event.", calendarEvents[1]),
			generateImage: false,
		},
		{
			key:           "calendar3",
			prompt:        fmt.Sprintf("You are a calendar assistant. The calendar event is below: %v.\n \n Write a very short comment on the calendar event.", calendarEvents[2]),
			generateImage: false,
		},
	}

	// concurrently generate updates for each prompt
	updates := make([]PromptResult, len(prompts))
	var wg sync.WaitGroup
	for i, promptValue := range prompts {
		wg.Add(1)
		go func(i int, promptValue prompt) {
			defer wg.Done()
			promptResult, err := o.generate(promptValue.prompt)
			if err != nil {
				writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
			}

			if promptValue.generateImage {
				imageURL, err := a.txt2img(promptResult)
				if err != nil {
					writeHttpError(w, http.StatusInternalServerError, "cannot get updates", err)
				}
				updates[i].ImageURL = imageURL
			}

			updates[i].Response = promptResult
			updates[i].Key = promptValue.key
		}(i, promptValue)
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
