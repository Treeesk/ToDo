package repos

// Тестирование аунтентификационных функций БД
import (
	"ProjectGo/backend/internal/config"
	"context"
	"testing"
	"time"
)

// Тест регистрации
func TestRegister(t *testing.T) {
	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	defer conn.Conn.Close()
	// Попытка создания существующего юзера
	_, _, err := conn.Register("Yar", "qwerty123", context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("expected error")
	}
}

// Тест логина
func TestLogin(t *testing.T) {
	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	defer conn.Conn.Close()
	// Верные данные
	id, err := conn.Login("Yar", "qwerty123", context.Background())
	if err != nil {
		t.Fatalf("error on login with valid data: %v", err)
	}
	// Проверка возврата верного id
	var trueId int
	if err = conn.Conn.QueryRow(context.Background(), "SELECT id FROM users WHERE user_login = $1", "Yar").Scan(&trueId); err != nil {
		t.Fatalf("unknown: %v", err)
	}
	if id != trueId {
		t.Fatal("The ID returned by Login does not match the one existing in the database")
	}
	// Неверный логин
	_, err = conn.Login("badLogin", "qwerty123", context.Background())
	if err == nil {
		t.Fatal("no error occurs when entering an invalid login")
	}
	// Неверный пароль
	_, err = conn.Login("Yar", "fdgfd", context.Background())
	if err == nil {
		t.Fatal("no error occurs when entering an invalid password")
	}
}
