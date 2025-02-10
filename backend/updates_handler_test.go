package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO:adjust tests so that prompt responses and image responses are mapped to the correct prompt
const (
	weatherPrompt  = "You are a weather assistant."
	newsPrompt     = "You are a news assistant."
	calendarPrompt = "You are a calendar assistant."
)

type mockOpenWebUIClient struct {
	generateCalls  int
	generateArgs   []string
	generateErrors []error
}

func (m *mockOpenWebUIClient) generate(prompt string) (string, error) {
	m.generateCalls++
	m.generateArgs = append(m.generateArgs, prompt)
	if len(m.generateErrors) > 0 {
		return "", m.generateErrors[m.generateCalls-1]
	}
	response := ""
	// return a response based on whether the prompt matches the beginning of the prompt
	if strings.HasPrefix(prompt, weatherPrompt) {
		response = "The weather is clear and sunny."
	} else if strings.HasPrefix(prompt, newsPrompt) {
		response = "The latest news is that the weather is clear and sunny."
	} else if strings.HasPrefix(prompt, calendarPrompt) {
		response = "The latest calendar event is that the weather is clear and sunny."
	}
	return response, nil
}

type mockAutomaticSDClient struct {
	txt2imgCalls  int
	txt2imgArgs   []string
	txt2imgErrors []error
}

func (m *mockAutomaticSDClient) txt2img(prompt string) (string, error) {
	m.txt2imgCalls++
	m.txt2imgArgs = append(m.txt2imgArgs, prompt)
	if len(m.txt2imgErrors) > 0 {
		return "", m.txt2imgErrors[m.txt2imgCalls-1]
	}
	return "test.com/image.jpg", nil
}

type mockWeatherClient struct {
	getCalls   int
	getReturns *weatherResult
	getErrors  []error
}

func (m *mockWeatherClient) get() (*weatherResult, error) {
	m.getCalls++
	if len(m.getErrors) > 0 {
		return nil, m.getErrors[m.getCalls-1]
	}
	return m.getReturns, nil
}

type mockNewsClient struct {
	getCalls   int
	getReturns []newsResult
	getErrors  []error
}

func (m *mockNewsClient) get() ([]newsResult, error) {
	m.getCalls++
	if len(m.getErrors) > 0 {
		return nil, m.getErrors[m.getCalls-1]
	}
	return m.getReturns, nil
}

type mockCalendarClient struct {
	getEventsCalls   int
	getEventsReturns []calendarEvent
	getEventsErrors  []error
}

func (m *mockCalendarClient) getEvents() ([]calendarEvent, error) {
	m.getEventsCalls++
	if len(m.getEventsErrors) > 0 {
		return nil, m.getEventsErrors[m.getEventsCalls-1]
	}
	return m.getEventsReturns, nil
}

func TestGetUpdatesSuccess(t *testing.T) {
	// set up test environment
	mockOpenWebUIClient := &mockOpenWebUIClient{
		generateErrors: []error{},
	}
	mockAutomaticSDClient := &mockAutomaticSDClient{
		txt2imgErrors: []error{},
	}
	mockWeatherClient := &mockWeatherClient{}
	mockNewsClient := &mockNewsClient{}
	mockCalendarClient := &mockCalendarClient{}

	// set up test data
	weatherValue := &weatherResult{
		Temp:    20.0,
		Weather: "Clear",
	}
	mockWeatherClient.getReturns = weatherValue
	mockWeatherClient.getErrors = []error{}

	newsValue := newsResult{
		Title:       "Test News",
		Description: "Test Description",
		URL:         "https://test.com",
	}
	mockNewsClient.getReturns = []newsResult{newsValue}
	mockNewsClient.getErrors = []error{}

	calendarEventValue := calendarEvent{
		Title:       "Test Calendar Event",
		Description: "Test Description",
		Start:       "2021-01-01",
		End:         "2021-01-01",
	}
	mockCalendarClient.getEventsReturns = []calendarEvent{calendarEventValue, calendarEventValue, calendarEventValue}
	mockCalendarClient.getEventsErrors = []error{}

	// set up test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// handle GET request for /updates
		if r.Method == "GET" {
			getUpdates(w, r, mockOpenWebUIClient, mockAutomaticSDClient, mockWeatherClient, mockNewsClient, mockCalendarClient)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}))
	defer server.Close()

	// set up test client
	client := server.Client()

	// set up test request
	req, err := http.NewRequest("GET", server.URL+"/updates", nil)
	require.NoError(t, err)

	// send request
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// check response status code
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// check response body
	var response []PromptResult
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	// check number of updates
	require.Len(t, response, 5)

	// check weather update
	require.Equal(t, "weather", response[0].Key)
	require.Equal(t, "The weather is clear and sunny.", response[0].Response)

	// check news update
	require.Equal(t, "news", response[1].Key)
	require.Equal(t, "The latest news is that the weather is clear and sunny.", response[1].Response)
	require.Equal(t, "test.com/image.jpg", response[1].ImageURL)

	// check calendar update
	require.Equal(t, "calendar1", response[2].Key)
	require.Equal(t, "The latest calendar event is that the weather is clear and sunny.", response[2].Response)
}
