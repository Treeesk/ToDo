package services

// Тестирование:
// верификации существующего и несуществующего токена, неверно указанный алгоритм шифрования, измененные поля токена

import (
	"ProjectGo/backend/internal/config"
	"ProjectGo/backend/internal/repos"
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")
	os.Exit(m.Run())
}

// Проверка создания и верификации токена
func TestCreateVerifyValid(t *testing.T) {
	cfg := config.Load()
	conn := repos.ConnUrlRepos(context.Background(), cfg)
	authService := NewAuthService(conn, cfg.JWTSecret)
	token, err := authService.CreateToken(1, time.Now().Add(time.Minute*5))
	if err != nil {
		t.Fatalf("Unknown error while creating the token: %v", err)
	}
	_, err = authService.VerifyToken(token)
	if err != nil {
		t.Fatalf("Error verifying token: %v", err)
	}
}

// Проверка ошибки верификации невалидного токена
func TestVerifyInvalid(t *testing.T) {
	cfg := config.Load()
	conn := repos.ConnUrlRepos(context.Background(), cfg)
	authService := NewAuthService(conn, cfg.JWTSecret)
	// Проверка на верификацию токена с неверным секретом
	claims := CustomClaims{
		1,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("apappa"))
	if err != nil {
		t.Fatalf("Error while creating the token: %v", err)
	}
	_, err = authService.VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("expected error")
	}

	// Проверка на неверно указанный метод шифрования
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["alg"] = "RSA256"
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Error while creating the token: %v", err)
	}
	_, err = authService.VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("signing method error expected")
	}

	// Проверка на ошибку верификации токена по истечению времени его существования
	claimsBadExp := CustomClaims{
		1,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		},
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claimsBadExp)
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Error while creating the token: %v", err)
	}
	_, err = authService.VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("expected error: out of time")
	}
}
