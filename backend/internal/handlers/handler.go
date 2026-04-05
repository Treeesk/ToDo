package handlers

import (
	"ProjectGo/backend/internal/services"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

type HandlerNotes struct {
	store *services.NotesStore
}

func NewHandlerNotes(store *services.NotesStore) *HandlerNotes {
	return &HandlerNotes{store: store}
}

type JsonError struct {
	Text string `json:"message"`
}

// Функция для отправки ошибки в JSON формате
func writeJsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(JsonError{Text: message})
}

// Функция возвращающая JSON с полным списком всех заметок
func (h *HandlerNotes) GetNotes(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(h.store.GetAll())
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "Error: processing in JSON format")
		log.Println("Error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Функция добавления заметки
func (h *HandlerNotes) AddNote(w http.ResponseWriter, r *http.Request) {
	type addn struct {
		Text string `json:"text"`
	}
	var note addn
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
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
		log.Println("Decode error: ", err)
		return
	}
	if strings.TrimSpace(note.Text) == "" {
		writeJsonError(w, http.StatusBadRequest, "text is required")
		log.Println("Error: text is required")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	h.store.Add(note.Text)
}
