package env

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type Env struct {
	ApiPort        string
	WeatherAPIKey  string
	WeatherBaseURL string
	APICEPBaseURL  string
}

func New(ctx context.Context) (*Env, error) {

	cfg := &Env{
		ApiPort:        getEnv("API_PORT"),
		WeatherAPIKey:  getEnv("WEATHER_API_KEY"),
		WeatherBaseURL: getEnv("WEATHER_BASE_URL"),
		APICEPBaseURL:  getEnv("API_CEP_BASE_URL"),
	}

	if cfg.ApiPort == "" {
		cfg.ApiPort = getEnv("PORT")
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func (e *Env) validate() error {
	missing := make([]string, 0)

	if e.ApiPort == "" {
		missing = append(missing, "API_PORT")
	}
	if e.WeatherAPIKey == "" {
		missing = append(missing, "WEATHER_API_KEY")
	}
	if e.WeatherBaseURL == "" {
		missing = append(missing, "WEATHER_BASE_URL")
	}
	if e.APICEPBaseURL == "" {
		missing = append(missing, "API_CEP_BASE_URL")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return nil
}
