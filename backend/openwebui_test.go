package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupOpenWebUIEnvVars(serverURL string) {
	// set up env vars
	err := os.Setenv("OPENWEBUI_BASE_URL", serverURL)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("OPENWEBUI_API_KEY", "test-api-key")
	if err != nil {
		panic(err)
	}

	err = os.Setenv("OPENWEBUI_MODEL_NAME", "test-model")
	if err != nil {
		panic(err)
	}
}

func setupOpenWebUIServer() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
	}))

	setupOpenWebUIEnvVars(mockServer.URL)

	return mockServer
}

func setupOpenWebUIServerWithInternalError() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "something bad happened"}`))
	}))

	setupOpenWebUIEnvVars(mockServer.URL)

	return mockServer
}

func setupOpenWebUIServerWithMalformedError() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`malformed json`))
	}))

	setupOpenWebUIEnvVars(mockServer.URL)

	return mockServer
}

func TestGenerateSuccess(t *testing.T) {
	// set up test environment
	setupOpenWebUIServer()

	// create a new openWebUI client
	openWebUIClient, err := newOpenWebUIClient()
	require.NoError(t, err)

	// generate some text
	response, err := openWebUIClient.generate("What's the weather like today?")
	require.NoError(t, err)

	// check that the response is not empty
	require.NotEmpty(t, response)
	require.Equal(t, "success", response)
}

func TestGenerateInternalError(t *testing.T) {
	// set up test environment
	setupOpenWebUIServerWithInternalError()

	// create a new openWebUI client
	openWebUIClient, err := newOpenWebUIClient()
	require.NoError(t, err)

	// generate some text
	response, err := openWebUIClient.generate("error")
	require.Error(t, err)
	require.Empty(t, response)
	require.Equal(t, "unexpected response code: 500", err.Error())
}

func TestGenerateMalformedResponse(t *testing.T) {
	// set up test environment
	setupOpenWebUIServerWithMalformedError()

	// create a new openWebUI client
	openWebUIClient, err := newOpenWebUIClient()
	require.NoError(t, err)

	// generate some text
	response, err := openWebUIClient.generate("error")
	require.Error(t, err)
	require.Empty(t, response)
	require.Contains(t, err.Error(), "cannot unmarshal response")
}
