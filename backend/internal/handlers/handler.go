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
// Получает JSON {"user_id": int}
// Возвращает JSON {"id": int, "user_id": int, "text": string}
func (h *HandlerNotes) GetNotes(w http.ResponseWriter, r *http.Request) {
	type user struct {
		User_id *int `json:"user_id"`
	}
	var us user
	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	if us.User_id == nil {
		writeJsonError(w, http.StatusBadRequest, "unknown user")
		log.Println("error: the user_id field is missing")
		return
	}
	var buf bytes.Buffer
	notes, err := h.store.GetAll(*(us.User_id))
	if err != nil {
		ErrorDB(w, err)
		log.Println("database error: ", err)
		return
	}
	err = json.NewEncoder(&buf).Encode(notes)
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
// Ожидается JSON вида {"user_id": int, "text": string}
func (h *HandlerNotes) AddNote(w http.ResponseWriter, r *http.Request) {
	type addn struct {
		User_id *int    `json:"user_id"`
		Text    *string `json:"text"`
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
	if note.User_id == nil {
		writeJsonError(w, http.StatusBadRequest, "unknown user")
		log.Println("error: the user_id field is missing")
		return
	}
	if strings.TrimSpace(*(note.Text)) == "" {
		writeJsonError(w, http.StatusBadRequest, "text of note is required")
		log.Println("error: text is required")
		return
	}
	err = h.store.Add(*(note.User_id), *(note.Text))
	if err != nil {
		ErrorDB(w, err)
		log.Println("database error: ", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Функция для удаления заметок
// Ожидается JSON вида {"user_id": int, "id": int}
func (h *HandlerNotes) DelNote(w http.ResponseWriter, r *http.Request) {
	type deln struct {
		User_id *int `json:"user_id"`
		ID      *int `json:"id"`
	}
	var note deln
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
	if note.User_id == nil {
		writeJsonError(w, http.StatusBadRequest, "the user_id field is missing")
		log.Println("error: the user_id field is missing")
		return
	}
	err = h.store.Del(*(note.User_id), *(note.ID))
	if err != nil {
		ErrorDB(w, err)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Функция редактирования заметок
// Ожидается JSON вида {"user_id:: int, "id": int, "text": string}
func (h *HandlerNotes) EditNote(w http.ResponseWriter, r *http.Request) {
	type editn struct {
		ID      *int    `json:"id"`
		User_id *int    `json:"user_id"`
		Text    *string `json:"text"`
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
	if note.User_id == nil {
		writeJsonError(w, http.StatusBadRequest, "the user_id field is missing")
		log.Println("error: the user_id field is missing")
		return
	}
	if strings.TrimSpace(*(note.Text)) == "" {
		writeJsonError(w, http.StatusBadRequest, "text is required")
		log.Println("error: text is required")
		return
	}
	err = h.store.Edit(*(note.User_id), *(note.ID), *(note.Text))
	if err != nil {
		ErrorDB(w, err)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
