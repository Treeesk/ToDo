package main

import (
	"ProjectGo/backend/internal/services"
	"ProjectGo/backend/internal/transport"
	"log"
	"net/http"
)

func main() {
	store := services.NewNotesStore() // инициализация заметок
	transport.Setuprouter(store)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
