package app

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetTemperature(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("cep")
	isValid := validateCep(zipCode)
	if !isValid {
		writePlainError(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	response := h.s.GetCurrentTemperatureByZipCode(zipCode)
	if response.Error != nil {
		writePlainError(w, response.Message, response.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var validateCepRegex = `^\d{8}$`

func validateCep(cep string) bool {
	match, err := regexp.MatchString(validateCepRegex, cep)
	if err != nil {
		return false
	}
	return match
}

func writePlainError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
