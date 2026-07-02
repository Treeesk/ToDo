package transport

import (
	"ProjectGo/backend/internal/handlers"
	"ProjectGo/backend/internal/services"
	"net/http"
)

func Setuprouter(mux *http.ServeMux, store *services.NotesStore, authService *services.AuthService) {
	handler := handlers.NewHandlerNotes(store, authService) // инициализация хендлер класса для работы с хендлер функциями заметок

	mux.HandleFunc("/api/", handler.GetNotes)
	mux.HandleFunc("/api/add/", handler.AddNote)
	mux.HandleFunc("/api/del/", handler.DelNote)
	mux.HandleFunc("/api/edit/", handler.EditNote)
	mux.HandleFunc("/api/register/", handler.Register)
	mux.HandleFunc("/api/login/", handler.Login)
	mux.HandleFunc("/api/logout/", handler.LogOut)
	mux.HandleFunc("/api/refresh/", handler.Refresh)
}
