package viacep

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Service struct {
	url     string
	timeout time.Duration
}

func NewService(url string) (*Service, error) {

	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}

	return &Service{
		url:     url,
		timeout: 5 * time.Second,
	}, nil
}

type Location struct {
	StatusCode int    `json:"-"`
	Message    string `json:"-"`
	Error      error  `json:"-"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

func (s *Service) GetLocationByZipCode(zipCode string) *Location {
	url := fmt.Sprintf("%s/ws/%s/json", s.url, zipCode)
	fmt.Println(url)
	httpClient := &http.Client{
		Timeout: s.timeout,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return &Location{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to get location by zip code",
			Error:      err,
		}
	}
	defer resp.Body.Close()

	var location *Location = &Location{
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		location.Error = fmt.Errorf("failed to get location by zip code: %s", resp.Status)
		if resp.StatusCode == http.StatusNotFound {
			location.Message = "can not find zipcode"
			location.Error = fmt.Errorf("can not find zipcode")
		}

		return location
	}

	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return &Location{
			StatusCode: http.StatusInternalServerError,
			Message:    "failed to decode response body",
			Error:      err,
		}
	}

	if strings.EqualFold("true", location.Erro) {
		location.StatusCode = http.StatusNotFound
		location.Message = "can not find zipcode"
		location.Error = fmt.Errorf("can not find zipcode")
	}

	return location
}
