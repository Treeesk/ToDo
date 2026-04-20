package main

import (
	"ProjectGo/backend/internal/config"
	"ProjectGo/backend/internal/repos"
	"ProjectGo/backend/internal/services"
	"ProjectGo/backend/internal/transport"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load() // конфиг с переменными окружения

	conn := repos.ConnUrlRepos(cfg)
	store := services.NewNotesStore() // инициализация заметок
	transport.Setuprouter(store)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, nil))
}
