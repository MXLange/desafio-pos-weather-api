package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	viacep "github.com/MXLange/desafio-pos-weather-api/internal/domain/via_cep"
	weatherapi "github.com/MXLange/desafio-pos-weather-api/internal/domain/weather_api"
)

func TestGetTemperatureSuccess(t *testing.T) {
	cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json" {
			t.Fatalf("unexpected ViaCEP path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"localidade":"Sao Paulo"}`))
	}))
	defer cepServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Fatalf("unexpected WeatherAPI path: %s", r.URL.Path)
		}
		if got, want := r.URL.Query().Get("q"), "Sao Paulo"; got != want {
			t.Fatalf("q = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer weatherServer.Close()

	handler := newTestHandler(t, cepServer.URL, weatherServer.URL)
	req := httptest.NewRequest(http.MethodGet, "/temperature?cep=01001000", nil)
	rec := httptest.NewRecorder()

	handler.GetTemperature(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%q", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body AppResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body.Temp_C != 28.5 || body.Temp_F != 83.3 || body.Temp_K != 301.65 {
		t.Fatalf("unexpected temperatures: %+v", body)
	}
}

func TestGetTemperatureInvalidZipcode(t *testing.T) {
	handler := newTestHandler(t, "http://example.com", "http://example.com")

	tests := []string{
		"/temperature",
		"/temperature?cep=1234567",
		"/temperature?cep=123456789",
		"/temperature?cep=12345-678",
		"/temperature?cep=abcdefgh",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt, nil)
			rec := httptest.NewRecorder()

			handler.GetTemperature(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
			}
			if got, want := rec.Body.String(), "invalid zipcode"; got != want {
				t.Fatalf("body = %q, want %q", got, want)
			}
		})
	}
}

func TestGetTemperatureZipcodeNotFound(t *testing.T) {
	cepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":"true"}`))
	}))
	defer cepServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("weather service should not be called when zipcode is not found")
	}))
	defer weatherServer.Close()

	handler := newTestHandler(t, cepServer.URL, weatherServer.URL)
	req := httptest.NewRequest(http.MethodGet, "/temperature?cep=99999999", nil)
	rec := httptest.NewRecorder()

	handler.GetTemperature(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if got, want := rec.Body.String(), "can not find zipcode"; got != want {
		t.Fatalf("body = %q, want %q", got, want)
	}
}

func newTestHandler(t *testing.T, cepBaseURL, weatherBaseURL string) *Handler {
	t.Helper()

	cepService, err := viacep.NewService(cepBaseURL)
	if err != nil {
		t.Fatalf("failed to create ViaCEP service: %v", err)
	}

	weatherService, err := weatherapi.NewService(weatherBaseURL, "test-key")
	if err != nil {
		t.Fatalf("failed to create WeatherAPI service: %v", err)
	}

	return NewHandler(NewService(cepService, weatherService))
}
