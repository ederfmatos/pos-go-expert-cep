package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	ViaCepResponse struct {
		Error      *string `json:"erro,omitempty"`
		Localidade string  `json:"localidade"`
	}

	ViaCepClient struct {
		baseURL string
		client  http.Client
	}
)

func NewViaCepClient(baseURL string) *ViaCepClient {
	return &ViaCepClient{
		baseURL: baseURL,
		client: http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (c ViaCepClient) GetAddress(postalCode string) (*ViaCepResponse, error) {
	url := fmt.Sprintf("%s/%s/json", c.baseURL, postalCode)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("send request: %v", err)
	}

	var response ViaCepResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode viacep response: %v", err)
	}

	_ = resp.Body.Close()
	return &response, nil
}
