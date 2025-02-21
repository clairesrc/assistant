package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type weatherClient interface {
	get() (*weatherResult, error)
}

type weather struct {
	apiKey    string
	latitude  string
	longitude string
	timezone  string
	baseURL   string

	value       weatherResult
	lastUpdated time.Time
}

type weatherResult struct {
	Temp    float64 `json:"temp"`
	Weather string  `json:"weather"`
}

const (
	weatherCacheDuration = 10 * time.Minute
)

func newWeatherClient() (weatherClient, error) {
	if os.Getenv("OPENWEATHER_API_KEY") == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY is not set")
	}

	if os.Getenv("OPENWEATHER_LATITUDE") == "" {
		return nil, fmt.Errorf("OPENWEATHER_LATITUDE is not set")
	}

	if os.Getenv("OPENWEATHER_LONGITUDE") == "" {
		return nil, fmt.Errorf("OPENWEATHER_LONGITUDE is not set")
	}

	if os.Getenv("OPENWEATHER_TIMEZONE") == "" {
		return nil, fmt.Errorf("OPENWEATHER_TIMEZONE is not set")
	}

	if os.Getenv("OPENWEATHER_BASE_URL") == "" {
		return nil, fmt.Errorf("OPENWEATHER_BASE_URL is not set")
	}

	return &weather{
		apiKey:    os.Getenv("OPENWEATHER_API_KEY"),
		latitude:  os.Getenv("OPENWEATHER_LATITUDE"),
		longitude: os.Getenv("OPENWEATHER_LONGITUDE"),
		timezone:  os.Getenv("OPENWEATHER_TIMEZONE"),
		baseURL:   os.Getenv("OPENWEATHER_BASE_URL"),
	}, nil
}

func (w *weather) get() (*weatherResult, error) {
	// if cache is fresh, return cached value
	if time.Since(w.lastUpdated) < weatherCacheDuration {
		return &weatherResult{
			Temp:    w.value.Temp,
			Weather: w.value.Weather,
		}, nil
	}

	// get weather data from openweathermap API
	url := fmt.Sprintf("%s/data/2.5/weather?lat=%s&lon=%s&appid=%s", w.baseURL, w.latitude, w.longitude, w.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get weather: %w", err)
	}

	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read weather: %w", err)
	}

	// parse weather data
	var weatherData struct {
		Current struct {
			Temperature float64 `json:"temp"`
			Weather     []struct {
				Main string `json:"main"`
			} `json:"weather"`
		} `json:"current"`
	}

	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return nil, fmt.Errorf("cannot parse weather: %w", err)
	}
	weather := weatherResult{
		Temp:    weatherData.Current.Temperature,
		Weather: weatherData.Current.Weather[0].Main,
	}

	w.value = weather
	w.lastUpdated = time.Now()
	return &weather, nil
}
