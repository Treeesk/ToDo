package main

import (
	"ProjectGo/backend/internal/config"
	"ProjectGo/backend/internal/repos"
	"ProjectGo/backend/internal/services"
	"ProjectGo/backend/internal/transport"
	"context"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg := config.Load() // конфиг с переменными окружения

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // контекст на подключение к Бд(защита от зависания)
	defer cancel()

	conn := repos.ConnUrlRepos(ctx, cfg)
	defer conn.Conn.Close()
	store := services.NewNotesStore(conn)
	authService := services.NewAuthService(cfg.JWTSecret)
	transport.Setuprouter(store, authService)
	log.Fatal(http.ListenAndServe(cfg.BaseURL, nil))
}
