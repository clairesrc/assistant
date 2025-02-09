package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/stretchr/testify/require"
)

func setupNewsClientEnvVars(serverURL string) {
	os.Setenv("NEWS_API_KEY", "test")
	os.Setenv("NEWS_BASE_URL", serverURL)
}

func setupNewsServerSuccess() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"title": "Test Title", "description": "Test Description", "url": "https://test.com"}]`))
	}))

	setupNewsClientEnvVars(server.URL)

	return server
}

func setupNewsServerInternalError() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	setupNewsClientEnvVars(server.URL)

	return server
}

func TestNewsClient_Get(t *testing.T) {
	server := setupNewsServerSuccess()
	defer server.Close()

	newsClient, err := newNewsClient()
	require.NoError(t, err)

	results, err := newsClient.get()
	require.NoError(t, err)
	require.Equal(t, 1, len(results))
	require.Equal(t, "Test Title", results[0].Title)
	require.Equal(t, "Test Description", results[0].Description)
	require.Equal(t, "https://test.com", results[0].URL)
}

func TestNewsClient_GetInternalError(t *testing.T) {
	server := setupNewsServerInternalError()
	defer server.Close()
	
	newsClient, err := newNewsClient()
	require.NoError(t, err)

	results, err := newsClient.get()
	require.Error(t, err)
	require.Nil(t, results)
}
