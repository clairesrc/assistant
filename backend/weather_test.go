package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/stretchr/testify/require"
)

const weatherSuccessResponse = `{
	"current": {
		"temp": 20,
		"weather": [
			{
				"main": "Clear"
			}
		]
	}
}`

func setupWeatherClientEnvVars(serverURL string) {
	os.Setenv("OPENWEATHER_API_KEY", "test")
	os.Setenv("OPENWEATHER_LATITUDE", "40.7128")
	os.Setenv("OPENWEATHER_LONGITUDE", "-74.0060")
	os.Setenv("OPENWEATHER_TIMEZONE", "America/Chicago")
	os.Setenv("OPENWEATHER_BASE_URL", serverURL)
}

func setupWeatherServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(weatherSuccessResponse))
	}))

	setupWeatherClientEnvVars(server.URL)

	return server
}

func setupWeatherServerWithInternalError() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	setupWeatherClientEnvVars(server.URL)

	return server
}

func TestWeatherClient_GetSuccess(t *testing.T) {
	server := setupWeatherServer()
	defer server.Close()

	weather, err := newWeatherClient()
	require.NoError(t, err)

	result, err := weather.get()
	require.NoError(t, err)
	require.Equal(t, 20, result.Temp)
	require.Equal(t, "Clear", result.Weather[0].Main)
}

func TestWeatherClient_GetInternalError(t *testing.T) {
	server := setupWeatherServerWithInternalError()
	defer server.Close()

	weather, err := newWeatherClient()
	require.NoError(t, err)

	result, err := weather.get()
	require.Error(t, err)
	require.Nil(t, result)
}