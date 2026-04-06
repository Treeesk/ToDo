package handlers

import (
	"ProjectGo/backend/internal/services"
	"bytes"
	"encoding/json"
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

// Функция возвращающая JSON с полным списком всех заметок
// Возвращается JSON {"id": id, "text": note}
func (h *HandlerNotes) GetNotes(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(h.store.GetAll())
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "Error: processing in JSON format")
		log.Println("error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Функция добавления заметки
// Ожидается JSON вида {"text": string}
func (h *HandlerNotes) AddNote(w http.ResponseWriter, r *http.Request) {
	type addn struct {
		Text *string `json:"text"`
	}
	var note addn
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	if note.Text == nil {
		writeJsonError(w, http.StatusBadRequest, "the text field is missing")
		log.Println("error: the text field is missing")
		return
	}
	if strings.TrimSpace(*(note.Text)) == "" {
		writeJsonError(w, http.StatusBadRequest, "text is required")
		log.Println("error: text is required")
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.store.Add(*(note.Text))
}

// Функция для удаления заметок
// Ожидается JSON вида {"id": int}
func (h *HandlerNotes) DelNote(w http.ResponseWriter, r *http.Request) {
	type deln struct {
		ID *int `json:"id"`
	}
	var id deln
	err := json.NewDecoder(r.Body).Decode(&id)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	if id.ID == nil {
		writeJsonError(w, http.StatusBadRequest, "the id field is missing")
		log.Println("error: the id field is missing")
		return
	}
	err = h.store.Del(*(id.ID))
	if err != nil {
		writeJsonError(w, 400, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Функция редактирования заметок
// Ожидается JSON вида {"id": int, "text": string}
func (h *HandlerNotes) EditNote(w http.ResponseWriter, r *http.Request) {
	type editn struct {
		ID   *int    `json:"id"`
		Text *string `json:"text"`
	}
	var note editn
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	if note.ID == nil {
		writeJsonError(w, http.StatusBadRequest, "the id field is missing")
		log.Println("error: the id field is missing")
		return
	}
	if note.Text == nil {
		writeJsonError(w, http.StatusBadRequest, "the text field is missing")
		log.Println("error: the text field is missing")
		return
	}
	if strings.TrimSpace(*(note.Text)) == "" {
		writeJsonError(w, http.StatusBadRequest, "text is required")
		log.Println("error: text is required")
		return
	}
	err = h.store.Edit(*(note.ID), *(note.Text))
	if err != nil {
		writeJsonError(w, 400, err.Error())
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
