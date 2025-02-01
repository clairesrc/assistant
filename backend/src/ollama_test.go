package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupEnvVars(serverURL string) {
	// set up env vars
	err := os.Setenv("OLLAMA_BASE_URL", serverURL)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("OLLAMA_API_KEY", "test-api-key")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("OLLAMA_MODEL_NAME", "test-model")
	if err != nil {
		panic(err)
	}
}

func setupMockServer() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
	}))

	setupEnvVars(mockServer.URL)

	return mockServer
}

func setupMockInternalErrorServer() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "something bad happened"}`))
	}))

	setupEnvVars(mockServer.URL)

	return mockServer
}

func setupMockMalformedErrorServer() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`malformed json`))
	}))

	setupEnvVars(mockServer.URL)

	return mockServer
}

func TestGenerateSuccess(t *testing.T) {
	// set up test environment
	setupMockServer()

	// create a new ollama client
	ollamaClient, err := newOllamaClient()
	require.NoError(t, err)

	// generate some text
	response, err := ollamaClient.generate("What's the weather like today?")
	require.NoError(t, err)

	// check that the response is not empty
	require.NotEmpty(t, response)
	require.Equal(t, "success", response)
}

func TestGenerateInternalError(t *testing.T) {
	// set up test environment
	setupMockInternalErrorServer()

	// create a new ollama client
	ollamaClient, err := newOllamaClient()
	require.NoError(t, err)

	// generate some text
	response, err := ollamaClient.generate("error")
	require.Error(t, err)
	require.Empty(t, response)
	require.Equal(t, "unexpected response code: 500", err.Error())
}

func TestGenerateMalformedResponse(t *testing.T) {
	// set up test environment
	setupMockMalformedErrorServer()

	// create a new ollama client
	ollamaClient, err := newOllamaClient()
	require.NoError(t, err)

	// generate some text
	response, err := ollamaClient.generate("error")
	require.Error(t, err)
	require.Empty(t, response)
	require.Contains(t, err.Error(), "cannot unmarshal response")
}
