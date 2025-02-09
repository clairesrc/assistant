package main

import (
	"fmt"
	"net/http"
	"os"
	"bytes"
	"io"
	"encoding/json"
)

type automaticSDClient struct {
	baseUrl   string
	modelName string
}

func newAutomaticSDClient() (*automaticSDClient, error) {
	if os.Getenv("AUTOMATIC1111_BASE_URL") == "" {
		return nil, fmt.Errorf("AUTOMATIC1111_BASE_URL env var is not set")
	}
	return &automaticSDClient{
		baseUrl:   os.Getenv("AUTOMATIC1111_BASE_URL"),
		modelName: os.Getenv("AUTOMATIC1111_MODEL_NAME"),
	}, nil
}

func (c *automaticSDClient) txt2img(prompt string) (string, error) {
	// payload for /sdapi/v1/txt2img endpoint
	payload := []byte(`{
		"prompt": "` + prompt + `",
		"negative_prompt": "",
		"width": 512,
		"height": 512,
		"enable_safety_checker": false,
		"enable_hr": true,
		"steps": 20,
		"hr_scale": 2,
		"hr_upscaler": "RealESRGAN_x4plus",
		"hr_second_pass_resolution": 1024,
		"hr_resize_mode": 1
	}`)

	// create request
	req, err := http.NewRequest("POST", c.baseUrl+"/sdapi/v1/txt2img", bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot send request: %w", err)
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read response: %w", err)
	}

	// parse response
	var response struct {
		Images []string `json:"images"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshal response: %w", err)
	}

	return response.Images[0], nil
}
