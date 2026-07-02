package handlers

import (
	"ProjectGo/backend/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
)

type HandlerNotes struct {
	store       *services.NotesStore
	authService *services.AuthService
}

func NewHandlerNotes(store *services.NotesStore, authService *services.AuthService) *HandlerNotes {
	return &HandlerNotes{
		store:       store,
		authService: authService,
	}
}

// Функция возвращающая JSON с полным списком всех заметок
// Возвращает JSON {"id": int, "user_id": int, "text": string}
func (h *HandlerNotes) GetNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: not allowed method")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()

	access_token, err := r.Cookie("access-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "missing cookie")
		log.Printf("unknown user: %v", err)
		return
	}
	user_id, err := h.authService.VerifyToken(access_token.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Printf("verify token fail: %v", err)
		return
	}
	var buf bytes.Buffer
	notes, err := h.store.GetAll(ctx, user_id)
	if err != nil {
		HandleError(w, err)
		log.Println("database error: ", err)
		return
	}
	err = json.NewEncoder(&buf).Encode(notes)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}

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
	if r.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
	}
	type addn struct {
		Text *string `json:"text"`
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()

	access_token, err := r.Cookie("access-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "missing cookie")
		log.Printf("unknown user: %v", err)
		return
	}
	user_id, err := h.authService.VerifyToken(access_token.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Printf("verify token fail: %v", err)
		return
	}

	var note addn
	err = json.NewDecoder(r.Body).Decode(&note)
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
		writeJsonError(w, http.StatusBadRequest, "text of note is required")
		log.Println("error: text is required")
		return
	}
	err = h.store.Add(ctx, user_id, *(note.Text))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}

		HandleError(w, err)
		log.Println("database error:", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Функция для удаления заметок
// Ожидается JSON вида {"id": int}
func (h *HandlerNotes) DelNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
	}
	type deln struct {
		ID *int `json:"id"`
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()

	access_token, err := r.Cookie("access-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "missing cookie")
		log.Printf("unknown user: %v", err)
		return
	}
	user_id, err := h.authService.VerifyToken(access_token.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Printf("verify token fail: %v", err)
		return
	}

	var note deln
	err = json.NewDecoder(r.Body).Decode(&note)
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
	err = h.store.Del(ctx, user_id, *(note.ID))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}

		HandleError(w, err)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Функция редактирования заметок
// Ожидается JSON вида {"id": int, "text": string}
func (h *HandlerNotes) EditNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
	}
	type editn struct {
		ID   *int    `json:"id"`
		Text *string `json:"text"`
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()

	access_token, err := r.Cookie("access-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "missing cookie")
		log.Printf("unknown user: %v", err)
		return
	}
	user_id, err := h.authService.VerifyToken(access_token.Value)
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Printf("verify token fail: %v", err)
		return
	}

	var note editn
	err = json.NewDecoder(r.Body).Decode(&note)
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
	err = h.store.Edit(ctx, user_id, *(note.ID), *(note.Text))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}

		HandleError(w, err)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
