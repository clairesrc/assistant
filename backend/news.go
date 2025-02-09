package main

import (
	"fmt"
	"net/http"
)

type news struct {
	apiKey string
	baseURL string
	
}

type newsResult struct {
	Title string `json:"title"`
	Description string `json:"description"`
	URL string `json:"url"`
}

func newNewsClient() (*news, error) {
	if os.Getenv("NEWS_API_KEY") == "" {
		return nil, fmt.Errorf("NEWS_API_KEY is not set")
	}

	if os.Getenv("NEWS_BASE_URL") == "" {
		return nil, fmt.Errorf("NEWS_BASE_URL is not set")
	}

	return &news{
		apiKey: os.Getenv("NEWS_API_KEY"),
		baseURL: os.Getenv("NEWS_BASE_URL"),
	}, nil
}

func (n *news) get() ([]newsResult, error) {
	url := fmt.Sprintf("%s/news?apiKey=%s", n.baseURL, n.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get news: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read news: %w", err)
	}
	
	var results []newsResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal news: %w", err)
	}

	return results, nil
}
