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
	if zipCode == "" {
		http.Error(w, "cep is required", http.StatusBadRequest)
		return
	}

	isValid := validateCep(zipCode)
	if !isValid {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	response := h.s.GetCurrentTemperatureByZipCode(zipCode)
	if response.Error != nil {
		http.Error(w, response.Message, response.StatusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

var validateCepRegex = `^\d{5}-?\d{3}$`

func validateCep(cep string) bool {
	match, err := regexp.MatchString(validateCepRegex, cep)
	if err != nil {
		return false
	}
	return match
}
