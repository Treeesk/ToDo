package transport

import (
	"ProjectGo/backend/internal/handlers"
	"ProjectGo/backend/internal/services"
	"net/http"
)

func Setuprouter(store *services.NotesStore, authService *services.AuthService) {
	handler := handlers.NewHandlerNotes(store, authService) // инициализация хендлер класса для работы с хендлер функциями заметок

	http.HandleFunc("/api/", handler.GetNotes)
	http.HandleFunc("/api/add/", handler.AddNote)
	http.HandleFunc("/api/del/", handler.DelNote)
	http.HandleFunc("/api/edit/", handler.EditNote)
	http.HandleFunc("/api/register/", handler.Register)
	http.HandleFunc("/api/login/", handler.Login)
	// http.HandleFunc("/api/register/", handler.Register) вызов функции из auth.go
}
