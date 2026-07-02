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

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux() // создание мультиплексора
	err := godotenv.Load()
	if err != nil {
		log.Println("config: error loading .env file, use base vars")
	}
	cfg := config.Load() // конфиг с переменными окружения
	server := &http.Server{
		Addr:              cfg.BaseURL,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // контекст на подключение к Бд(защита от зависания)
	defer cancel()

	conn := repos.ConnUrlRepos(ctx, cfg)
	defer conn.Conn.Close()
	store := services.NewNotesStore(conn)
	authService := services.NewAuthService(conn, cfg.JWTSecret)
	transport.Setuprouter(mux, store, authService)
	log.Fatal(server.ListenAndServe())
}
