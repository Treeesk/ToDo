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
	id, refresh_token, err := conn.Login("Yar", "qwerty123", context.Background(), time.Now().AddDate(0, 0, 30))
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
	err = conn.LogOut(context.Background(), refresh_token)
	if err != nil {
		t.Fatalf("failed logout: %v", err)
	}
	// Неверный логин
	_, _, err = conn.Login("badLogin", "qwerty123", context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("no error occurs when entering an invalid login")
	}
	// Неверный пароль
	_, _, err = conn.Login("Yar", "fdgfd", context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("no error occurs when entering an invalid password")
	}
	// Попытка залогиниться более чем 100 раз
	for i := 0; i < 101; i++ {
		id, _, err = conn.Login("Yar", "qwerty123", context.Background(), time.Now().AddDate(0, 0, 30))
		if err != nil {
			t.Fatalf("error on login with valid data: %v", err)
		}
	}
	tag, err := conn.Conn.Exec(context.Background(), "SELECT * FROM refresh_tokens WHERE user_id = $1", id)
	if tag.RowsAffected() > 100 {
		t.Fatal("error: more refresh tokens have been created than allowed")
	}
	_, err = conn.Conn.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE user_id = $1", id)
}

// Тест выхода из профиля
func TestLogOut(t *testing.T) {
	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	defer conn.Conn.Close()

	// логин и логаут
	id, refresh_token, err := conn.Login("Test", "test123", context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	err = conn.LogOut(context.Background(), refresh_token)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}
	tag, err := conn.Conn.Exec(context.Background(), "SELECT * FROM refresh_tokens WHERE user_id = $1", id)
	if err != nil {
		t.Fatalf("error accessing the database: %v", err)
	}
	if tag.RowsAffected() != 0 {
		t.Fatal("error: refresh token was not deleted")
	}
}

// Тест refresh token
func TestRefreshingToken(t *testing.T) {
	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	defer conn.Conn.Close()

	// Проверка создания нового рефреш токена
	_, refresh_token, err := conn.Login("Test", "test123", context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	_, new_refresh_token, err := conn.Refresh(refresh_token, time.Now().Add(5*time.Minute), context.Background())
	if err != nil {
		t.Fatalf("unknown error on Refresh: %v", err)
	}
	if refresh_token == new_refresh_token {
		t.Fatal("error: old refresh token equels new refresh token")
	}
	err = conn.LogOut(context.Background(), new_refresh_token)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}

	// Проверка ошибки времени существования токена
	_, refresh_token, err = conn.Login("Test", "test123", context.Background(), time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	_, _, err = conn.Refresh(refresh_token, time.Now().Add(5*time.Minute), context.Background())
	if err == nil {
		t.Fatal("Expired token error missed")
	}
	err = conn.LogOut(context.Background(), refresh_token)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}

	// Попытка зарефрешить неизвестный токен
	_, _, err = conn.Refresh(refresh_token, time.Now().Add(5*time.Minute), context.Background())
	if err == nil {
		t.Fatal("error missed: refresh of an unknown token")
	}
}
