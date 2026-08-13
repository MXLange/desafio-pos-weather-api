package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	url     string
	key     string
	timeout time.Duration
}

func NewService(url, key string) (*Service, error) {
	if url == "" {
		return nil, fmt.Errorf("url cannot be empty")
	}
	if key == "" {
		return nil, fmt.Errorf("key cannot be empty")
	}

	return &Service{
		url:     url,
		key:     key,
		timeout: 5 * time.Second,
	}, nil
}

type WeatherResponse struct {
	StatusCode int     `json:"-"`
	Message    string  `json:"-"`
	Error      error   `json:"-"`
	Current    Current `json:"current"`
}

type Current struct {
	TempC float64 `json:"temp_c"`
}

func (s *Service) GetCurrentCelciusByCity(city string) *WeatherResponse {
	url := fmt.Sprintf("%s/v1/current.json?key=%s&q=%s", s.url, s.key, city)

	httpClient := &http.Client{
		Timeout: s.timeout,
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return &WeatherResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch weather data",
			Error:      err,
		}
	}
	defer resp.Body.Close()

	var weatherResp *WeatherResponse = &WeatherResponse{
		StatusCode: resp.StatusCode,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		weatherResp.Error = fmt.Errorf("failed to get weather data: %s", resp.Status)
		if resp.StatusCode == http.StatusNotFound {
			weatherResp.Message = "can not find city"
			weatherResp.Error = fmt.Errorf("can not find city")
		}

		return weatherResp
	}

	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return &WeatherResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to decode weather data",
			Error:      err,
		}
	}

	return weatherResp
}
