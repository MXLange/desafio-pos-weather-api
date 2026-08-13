package helpers

func CelciusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

func CelciusToKelvin(c float64) float64 {
	return c + 273.15
}
