package repos

// Тестирование аунтентификационных функций БД
import (
	"ProjectGo/backend/internal/config"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

const testPassword = "qwerty123"

func openTestConn(t *testing.T) *ConnRepo {
	t.Helper()

	cfg := config.Load()
	conn := ConnUrlRepos(context.Background(), cfg)
	t.Cleanup(func() {
		conn.Conn.Close()
	})

	return conn
}

func createTestUser(t *testing.T, conn *ConnRepo) string {
	t.Helper()

	login := uniqueTestLogin(t)
	_, refreshToken, err := conn.Register(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("failed to register test user: %v", err)
	}
	if err := conn.LogOut(context.Background(), refreshToken); err != nil {
		t.Fatalf("failed to clean registration refresh token: %v", err)
	}

	t.Cleanup(func() {
		_, _ = conn.Conn.Exec(context.Background(), "DELETE FROM users WHERE user_login = $1", login)
	})

	return login
}

func uniqueTestLogin(t *testing.T) string {
	t.Helper()

	name := strings.NewReplacer("/", "_", " ", "_").Replace(t.Name())
	return fmt.Sprintf("test_%s_%d", name, time.Now().UnixNano())
}

// Тест регистрации
func TestRegister(t *testing.T) {
	conn := openTestConn(t)
	login := uniqueTestLogin(t)
	t.Cleanup(func() {
		_, _ = conn.Conn.Exec(context.Background(), "DELETE FROM users WHERE user_login = $1", login)
	})

	id, refreshToken, err := conn.Register(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("error on register with valid data: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive user id, got %d", id)
	}
	if refreshToken == "" {
		t.Fatal("expected refresh token")
	}

	_, _, err = conn.Register(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("expected duplicate login error")
	}
}

// Тест логина
func TestLogin(t *testing.T) {
	conn := openTestConn(t)
	login := createTestUser(t, conn)

	// Верные данные
	id, refreshToken, err := conn.Login(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("error on login with valid data: %v", err)
	}
	// Проверка возврата верного id
	var trueId int
	if err = conn.Conn.QueryRow(context.Background(), "SELECT id FROM users WHERE user_login = $1", login).Scan(&trueId); err != nil {
		t.Fatalf("unknown: %v", err)
	}
	if id != trueId {
		t.Fatal("The ID returned by Login does not match the one existing in the database")
	}
	err = conn.LogOut(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("failed logout: %v", err)
	}
	// Неверный логин
	_, _, err = conn.Login(uniqueTestLogin(t), testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("no error occurs when entering an invalid login")
	}
	// Неверный пароль
	_, _, err = conn.Login(login, "fdgfd", context.Background(), time.Now().AddDate(0, 0, 30))
	if err == nil {
		t.Fatal("no error occurs when entering an invalid password")
	}
	// Попытка залогиниться более чем 100 раз
	for i := 0; i < 101; i++ {
		id, _, err = conn.Login(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
		if err != nil {
			t.Fatalf("error on login with valid data: %v", err)
		}
	}
	var tokenCount int
	err = conn.Conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1", id).Scan(&tokenCount)
	if err != nil {
		t.Fatalf("error counting refresh tokens: %v", err)
	}
	if tokenCount > 100 {
		t.Fatal("error: more refresh tokens have been created than allowed")
	}
}

// Тест выхода из профиля
func TestLogOut(t *testing.T) {
	conn := openTestConn(t)
	login := createTestUser(t, conn)

	// логин и логаут
	id, refreshToken, err := conn.Login(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	err = conn.LogOut(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}
	var tokenCount int
	err = conn.Conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1", id).Scan(&tokenCount)
	if err != nil {
		t.Fatalf("error accessing the database: %v", err)
	}
	if tokenCount != 0 {
		t.Fatal("error: refresh token was not deleted")
	}
}

// Тест refresh token
func TestRefreshingToken(t *testing.T) {
	conn := openTestConn(t)
	login := createTestUser(t, conn)

	// Проверка создания нового рефреш токена
	_, refreshToken, err := conn.Login(login, testPassword, context.Background(), time.Now().AddDate(0, 0, 30))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	_, newRefreshToken, err := conn.Refresh(refreshToken, time.Now().Add(5*time.Minute), context.Background())
	if err != nil {
		t.Fatalf("unknown error on Refresh: %v", err)
	}
	if refreshToken == newRefreshToken {
		t.Fatal("error: old refresh token equels new refresh token")
	}
	err = conn.LogOut(context.Background(), newRefreshToken)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}

	// Проверка ошибки времени существования токена
	_, refreshToken, err = conn.Login(login, testPassword, context.Background(), time.Now().Add(-time.Minute))
	if err != nil {
		t.Fatalf("unknown error on Login: %v", err)
	}
	_, _, err = conn.Refresh(refreshToken, time.Now().Add(5*time.Minute), context.Background())
	if err == nil {
		t.Fatal("Expired token error missed")
	}
	err = conn.LogOut(context.Background(), refreshToken)
	if err != nil {
		t.Fatalf("error on LogOut: %v", err)
	}

	// Попытка зарефрешить неизвестный токен
	_, _, err = conn.Refresh(refreshToken, time.Now().Add(5*time.Minute), context.Background())
	if err == nil {
		t.Fatal("error missed: refresh of an unknown token")
	}
}
