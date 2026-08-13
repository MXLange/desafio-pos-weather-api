package env

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewReadsEnvironmentWithoutDotEnvFile(t *testing.T) {
	t.Setenv("API_PORT", "8080")
	t.Setenv("WEATHER_API_KEY", "test-key")
	t.Setenv("WEATHER_BASE_URL", "https://api.weatherapi.com")
	t.Setenv("API_CEP_BASE_URL", "https://viacep.com.br")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(cwd)

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	cfg, err := New(context.Background())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if got, want := cfg.ApiPort, "8080"; got != want {
		t.Fatalf("ApiPort = %q, want %q", got, want)
	}
	if got, want := cfg.WeatherAPIKey, "test-key"; got != want {
		t.Fatalf("WeatherAPIKey = %q, want %q", got, want)
	}
}

func TestNewUsesCloudRunPortWhenAPIPortIsMissing(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("WEATHER_API_KEY", "test-key")
	t.Setenv("WEATHER_BASE_URL", "https://api.weatherapi.com")
	t.Setenv("API_CEP_BASE_URL", "https://viacep.com.br")

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(cwd)

	tmp := t.TempDir()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	cfg, err := New(context.Background())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if got, want := cfg.ApiPort, "9090"; got != want {
		t.Fatalf("ApiPort = %q, want %q", got, want)
	}
}

func TestNewReadsDotEnvFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	defer os.Chdir(cwd)

	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, ".env"), []byte(`API_PORT=8081
WEATHER_API_KEY=file-key
WEATHER_BASE_URL=https://api.weatherapi.com
API_CEP_BASE_URL=https://viacep.com.br
`), 0600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	cfg, err := New(context.Background())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if got, want := cfg.ApiPort, "8081"; got != want {
		t.Fatalf("ApiPort = %q, want %q", got, want)
	}
	if got, want := cfg.WeatherAPIKey, "file-key"; got != want {
		t.Fatalf("WeatherAPIKey = %q, want %q", got, want)
	}
}
