package handlers

import (
	"ProjectGo/backend/internal/customerrors"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
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
	case errors.Is(err, context.Canceled):
		return
	case errors.Is(err, context.DeadlineExceeded):
		writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
		log.Println(err)
		return
	case errors.Is(err, customerrors.ErrTokenCreate):
		writeJsonError(w, http.StatusInternalServerError, "server error")
		log.Println(err)
		return
	case errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrPasswordTooLong):
		writeJsonError(w, http.StatusBadRequest, "too long or too short password")
		log.Println(err)
		return
	case errors.As(err, &pgerr):
		switch pgerr.Code {
		case "23503":
			writeJsonError(w, http.StatusNotFound, "Error: foreign_key_violation")
		case "22P02":
			writeJsonError(w, http.StatusBadRequest, "Error: invalid data type")
		case "23505":
			writeJsonError(w, http.StatusBadRequest, "Error: A user with that login already exists.")
		case "22001":
			writeJsonError(w, http.StatusBadRequest, "Error: input is too long")
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
