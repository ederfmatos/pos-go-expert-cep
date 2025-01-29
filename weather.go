package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	WeatherResponse struct {
		Current WeatherCurrent `json:"current"`
	}

	WeatherCurrent struct {
		TempC float64 `json:"temp_c"`
	}

	WeatherClient struct {
		baseURL string
		apiKey  string
		client  http.Client
	}
)

func NewWeatherClient(baseURL, apiKey string) *WeatherClient {
	return &WeatherClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c WeatherClient) GetWeather(city string) (*WeatherResponse, error) {
	weatherAPIURL := fmt.Sprintf("%s?key=%s&q=%s", c.baseURL, c.apiKey, url.QueryEscape(city))

	resp, err := c.client.Get(weatherAPIURL)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}

	var response WeatherResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode weather response: %v", err)
	}

	_ = resp.Body.Close()
	return &response, nil
}
