package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/stretchr/testify/require"
)

func setupAutomaticSDEnvVars(serverURL string) {
	// set up env vars
	err := os.Setenv("AUTOMATIC1111_BASE_URL", serverURL)
	if err != nil {
		panic(err)
	}

	err = os.Setenv("AUTOMATIC1111_MODEL_NAME", "test-model")
	if err != nil {
		panic(err)
	}
}

func setupAutomaticSDServer() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"images": ["test-image-url"]}`))
	}))

	setupAutomaticSDEnvVars(mockServer.URL)

	return mockServer
}

func setupAutomaticSDServerWithInternalError() *httptest.Server {
	// set up mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	setupAutomaticSDEnvVars(mockServer.URL)

	return mockServer
}

func TestAutomaticSDClient_Txt2ImgSuccess(t *testing.T) {
	// set up test environment
	setupAutomaticSDServer()

	// create a new automaticSD client
	automaticSDClient, err := newAutomaticSDClient()
	require.NoError(t, err)

	// call the txt2img method
	imageURL, err := automaticSDClient.txt2img("test prompt")
	require.NoError(t, err)
	require.Equal(t, "test-image-url", imageURL)
}

func TestAutomaticSDClient_Txt2ImgInternalError(t *testing.T) {
	// set up test environment
	setupAutomaticSDServerWithInternalError()

	// create a new automaticSD client
	automaticSDClient, err := newAutomaticSDClient()
	require.NoError(t, err)

	// call the txt2img method
	imageURL, err := automaticSDClient.txt2img("test prompt")
	require.Error(t, err)
	require.Empty(t, imageURL)
}
