package handlers

import (
	"ProjectGo/backend/internal/services"
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

	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
