package repos

import (
	"ProjectGo/backend/internal/config"
	"context"
	"testing"
	"time"
)

func TestCancel(t *testing.T) {
	// создание контекста, передача его и закрытие для проверки, что контекст прекращает обращение к бд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	cfg := &config.Config{
		BaseURL:    "localhost:8080",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "myuser",
		DBPassword: "IlikeCoding666",
		DBName:     "noteapp",
	}
	conn := ConnUrlRepos(cfg)
	defer conn.Conn.Close()

	_, err := conn.Conn.Query(ctx, "SELECT pg_sleep(1)")
	if err == nil {
		t.Fatal("expected error")
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Fatalf("expected deadline exceeded, got err=%v ctxErr=%v", err, ctx.Err())
	}
}
