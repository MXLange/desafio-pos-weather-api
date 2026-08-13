package weatherapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCurrentCelciusByCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Fatalf("path = %q, want %q", r.URL.Path, "/v1/current.json")
		}
		if got, want := r.URL.Query().Get("key"), "test-key"; got != want {
			t.Fatalf("key = %q, want %q", got, want)
		}
		if got, want := r.URL.Query().Get("q"), "Sao Paulo"; got != want {
			t.Fatalf("q = %q, want %q", got, want)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	service, err := NewService(server.URL, "test-key")
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}

	weather := service.GetCurrentCelciusByCity("Sao Paulo")
	if weather.Error != nil {
		t.Fatalf("GetCurrentCelciusByCity returned error: %v", weather.Error)
	}
	if got, want := weather.Current.TempC, 28.5; got != want {
		t.Fatalf("TempC = %v, want %v", got, want)
	}
}
