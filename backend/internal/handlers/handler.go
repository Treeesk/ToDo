package handlers

import (
	"ProjectGo/backend/internal/services"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type HandlerNotes struct {
	store *services.NotesStore
}

func NewHandlerNotes(store *services.NotesStore) *HandlerNotes {
	return &HandlerNotes{store: store}
}

// Функция возвращающая JSON с полным списком всех заметок
func (h *HandlerNotes) GetNotes(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(h.store.GetAll())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
