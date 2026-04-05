package transport

import (
	"ProjectGo/backend/internal/handlers"
	"ProjectGo/backend/internal/services"
	"net/http"
)

func Setuprouter(store *services.NotesStore) {
	handler := handlers.NewHandlerNotes(store) // инициализация хендлер класса для работы с хендлер функциями заметок

	http.HandleFunc("/api/", handler.GetNotes)
	http.HandleFunc("/api/add/", handler.AddNote)
}
