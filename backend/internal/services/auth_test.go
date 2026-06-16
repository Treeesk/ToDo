package services

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Тестирование:
// верификации существующего и несуществующего токена, неверно указанный алгоритм шифрования, измененные поля токена

// Проверка создания и верификации токена
func TestCreateVerifyValid(t *testing.T) {
	authService := NewAuthService("secret")
	token, err := authService.CreateToken(1)
	if err != nil {
		t.Fatalf("Unknown error while creating the token: %v", err)
	}
	err = authService.VerifyToken(token)
	if err != nil {
		t.Fatalf("Error verifying token: %v", err)
	}
}

// Проверка ошибки верификации невалидного токена
func TestVerifyInvalid(t *testing.T) {
	authService := NewAuthService("secret")
	// Проверка на верификацию токена с неверным секретом
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": 1,
			"iat": jwt.NewNumericDate(time.Now()),
			"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		})
	tokenString, err := token.SignedString([]byte("apappa"))
	if err != nil {
		t.Fatalf("Error while creating the token: %v", err)
	}
	err = authService.VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("expected error")
	}
	// Проверка на неверно указанный метод шифрования
	token = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": 1,
			"iat": jwt.NewNumericDate(time.Now()),
			"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		})
	token.Header["alg"] = "RSA256"
	tokenString, err = token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("Error while creating the token: %v", err)
	}
	err = authService.VerifyToken(tokenString)
	if err == nil {
		t.Fatalf("signing method error expected")
	}
}
