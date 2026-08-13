package helpers

import "testing"

func TestTemperatureConversions(t *testing.T) {
	celsius := 28.5

	if got, want := CelciusToFahrenheit(celsius), 83.3; got != want {
		t.Fatalf("CelciusToFahrenheit(%v) = %v, want %v", celsius, got, want)
	}

	if got, want := CelciusToKelvin(celsius), 301.65; got != want {
		t.Fatalf("CelciusToKelvin(%v) = %v, want %v", celsius, got, want)
	}
}
