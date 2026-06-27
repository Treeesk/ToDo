package handlers

import (
	"ProjectGo/backend/internal/customerrors"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
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

// Функция отправки ошибки во взаимодействии с БД в Response в формате JSON
func HandleError(w http.ResponseWriter, err error) {
	var pgerr *pgconn.PgError
	var notfounderr *customerrors.ErrorNotFound
	var userErr *customerrors.UserError
	switch {
	case errors.As(err, &pgerr):
		switch pgerr.Code {
		case "23503":
			writeJsonError(w, http.StatusNotFound, "Error: foreign_key_violation")
		case "22P02":
			writeJsonError(w, http.StatusBadRequest, "Error: invalid data type")
		case "23505":
			writeJsonError(w, http.StatusBadRequest, "Error: A user with that login already exists.")
		case "22001":
			writeJsonError(w, http.StatusBadRequest, "Error: data field missing")
		case "08006":
			writeJsonError(w, http.StatusServiceUnavailable, "Error: connection DB failure")
		case "22021":
			writeJsonError(w, http.StatusBadRequest, "Error: invalid byte sequence")
		}
	case errors.As(err, &notfounderr):
		writeJsonError(w, http.StatusNotFound, notfounderr.What)
	case errors.As(err, &userErr):
		writeJsonError(w, http.StatusUnauthorized, userErr.What)
	default:
		writeJsonError(w, http.StatusInternalServerError, "DB error")
	}
}
