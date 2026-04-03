package handlers

import (
	"ProjectGo/backend/internal/services"
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
	w.Header().Set("Content-Type", "application/json")
	encod := json.NewEncoder(w)
	err := encod.Encode(h.store.GetAll())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}
}
