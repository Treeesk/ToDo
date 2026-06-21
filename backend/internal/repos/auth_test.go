package repos

// Тестирование аунтентификационных функций БД
import (
	"ProjectGo/backend/internal/config"
	"context"
	"testing"
)

// Тест регистрации
func TestRegister(t *testing.T) {
	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	defer conn.Conn.Close()
	// Попытка создания существующего юзера
	_, err := conn.Register("Yar", "qwerty123", context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}
