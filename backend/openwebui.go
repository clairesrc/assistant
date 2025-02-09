package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type openWebUIClient interface {
	generate(prompt string) (string, error)
}

type openWebUI struct {
	baseUrl   string
	apiKey    string
	modelName string
}

func newOpenWebUIClient() (openWebUIClient, error) {
	if os.Getenv("OPENWEBUI_BASE_URL") == "" {
		return nil, fmt.Errorf("OPENWEBUI_BASE_URL env var is not set")
	}
	if os.Getenv("OPENWEBUI_API_KEY") == "" {
		return nil, fmt.Errorf("OPENWEBUI_API_KEY env var is not set")
	}
	if os.Getenv("OPENWEBUI_MODEL_NAME") == "" {
		return nil, fmt.Errorf("OPENWEBUI_MODEL_NAME env var is not set")
	}
	return &openWebUI{
		baseUrl:   os.Getenv("OPENWEBUI_BASE_URL"),
		apiKey:    os.Getenv("OPENWEBUI_API_KEY"),
		modelName: os.Getenv("OPENWEBUI_MODEL_NAME"),
	}, nil
}

func (o *openWebUI) generate(prompt string) (string, error) {
	// payload for /api/generate endpoint
	updatesPayload := []byte(`{
		"model": "` + o.modelName + `",
		"messages": [
			{
				"role": "user",
				"content": "` + prompt + `"
			}
		]
		}`)

	// create request
	req, err := http.NewRequest("POST", o.baseUrl+"/api/chat/completions", bytes.NewBuffer(updatesPayload))
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot send request: %w", err)
	}
	defer resp.Body.Close()

	// validate response code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	// read response
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBody := buf.String()

	// parse response
	var response struct {
		Response string `json:"response"`
	}
	err = json.Unmarshal([]byte(respBody), &response)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal response: %w", err)
	}

	return response.Response, nil
}
