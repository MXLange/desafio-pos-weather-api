package viacep

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocationByZipCodeSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ws/01001000/json" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/ws/01001000/json")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"localidade":"Sao Paulo"}`))
	}))
	defer server.Close()

	service, err := NewService(server.URL)
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}

	location := service.GetLocationByZipCode("01001000")
	if location.Error != nil {
		t.Fatalf("GetLocationByZipCode returned error: %v", location.Error)
	}
	if got, want := location.Localidade, "Sao Paulo"; got != want {
		t.Fatalf("Localidade = %q, want %q", got, want)
	}
}

func TestGetLocationByZipCodeNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	service, err := NewService(server.URL)
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}

	location := service.GetLocationByZipCode("99999999")
	if location.Error == nil {
		t.Fatal("expected zipcode not found error")
	}
	if got, want := location.StatusCode, http.StatusNotFound; got != want {
		t.Fatalf("StatusCode = %d, want %d", got, want)
	}
	if got, want := location.Message, "can not find zipcode"; got != want {
		t.Fatalf("Message = %q, want %q", got, want)
	}
}
