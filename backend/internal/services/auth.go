package services

// Создание и проверка JWT токенов. Регистрация и логин пользователей.

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

type CustomClaims struct {
	User_id int `json:"user_id"`
	jwt.RegisteredClaims
}

// Функция по созданию JWT
func (auth *AuthService) CreateToken(user_id int) (string, error) {
	claims := CustomClaims{
		user_id,
		jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(auth.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Функция проверки JWT
func (auth *AuthService) VerifyToken(tokenString string) error {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method") // явная документация ожидаемого алгоритма
		}
		return []byte(auth.jwtSecret), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}
