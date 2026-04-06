package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonError struct {
	Text string `json:"message"`
}

// Функция для отправки ошибки в JSON формате
func writeJsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonError{Text: message})
}

// Функция отправки ошибки в Response в формате JSON
func jsonDecodeError(w http.ResponseWriter, err error) {
	var syntxerr *json.SyntaxError
	var typeerr *json.UnmarshalTypeError
	switch {
	case errors.Is(err, io.EOF):
		writeJsonError(w, http.StatusBadRequest, "Error: body of json empty")
	case errors.As(err, &syntxerr):
		writeJsonError(w, http.StatusBadRequest, "Error: syntax error")
	case errors.As(err, &typeerr):
		writeJsonError(w, http.StatusBadRequest, "Error: JSON value is not appropriate for a given target type, or if a JSON number overflows the target type")
	default:
		writeJsonError(w, http.StatusBadRequest, "Error: invalid request body")
	}
}
