package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/stretchr/testify/require"
)

func setupCalendarClientEnvVars(serverURL string) {
	os.Setenv("CALENDAR_API_KEY", "test")
	os.Setenv("CALENDAR_BASE_URL", serverURL)
}

func setupCalendarServerSuccess() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"title": "Test Title", "description": "Test Description", "start": "2024-01-01", "end": "2024-01-02"}]`))
	}))

	setupCalendarClientEnvVars(server.URL)

	return server
}

func setupCalendarServerInternalError() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	
	setupCalendarClientEnvVars(server.URL)

	return server
}

func TestCalendarClient_GetEvents(t *testing.T) {
	server := setupCalendarServerSuccess()
	defer server.Close()

	calendarClient, err := newCalendarClient()
	require.NoError(t, err)
	
	events, err := calendarClient.getEvents()
	require.NoError(t, err)
	require.Equal(t, 1, len(events))
	require.Equal(t, "Test Title", events[0].Title)
	require.Equal(t, "Test Description", events[0].Description)
	require.Equal(t, "2024-01-01", events[0].Start)
	require.Equal(t, "2024-01-02", events[0].End)
}

func TestCalendarClient_GetEventsInternalError(t *testing.T) {
	server := setupCalendarServerInternalError()
	defer server.Close()

	calendarClient, err := newCalendarClient()
	require.NoError(t, err)

	events, err := calendarClient.getEvents()
	require.Error(t, err)
	require.Nil(t, events)
}

