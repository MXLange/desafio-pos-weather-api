package app

import (
	viacep "github.com/MXLange/desafio-pos-weather-api/internal/domain/via_cep"
	weatherapi "github.com/MXLange/desafio-pos-weather-api/internal/domain/weather_api"
	"github.com/MXLange/desafio-pos-weather-api/internal/helpers"
)

type Service struct {
	vc *viacep.Service
	wa *weatherapi.Service
}

func NewService(vc *viacep.Service, wa *weatherapi.Service) *Service {
	return &Service{
		vc: vc,
		wa: wa,
	}
}

type AppResponse struct {
	StatusCode int     `json:"-"`
	Error      error   `json:"-"`
	Message    string  `json:"message,omitempty"`
	Temp_C     float64 `json:"temp_C"`
	Temp_F     float64 `json:"temp_F"`
	Temp_K     float64 `json:"temp_K"`
}

func (s *Service) GetCurrentTemperatureByZipCode(zipCode string) *AppResponse {
	location := s.vc.GetLocationByZipCode(zipCode)
	if location.Error != nil {
		return &AppResponse{
			StatusCode: location.StatusCode,
			Message:    location.Message,
			Error:      location.Error,
		}
	}

	weather := s.wa.GetCurrentCelciusByCity(location.Localidade)
	if weather.Error != nil {
		return &AppResponse{
			StatusCode: weather.StatusCode,
			Message:    weather.Message,
			Error:      weather.Error,
		}
	}

	return &AppResponse{
		Temp_C: weather.Current.TempC,
		Temp_F: helpers.CelciusToFahrenheit(weather.Current.TempC),
		Temp_K: helpers.CelciusToKelvin(weather.Current.TempC),
	}
}
