package services

// Создание и проверка JWT токенов. Регистрация и логин пользователей.

import (
	"ProjectGo/backend/internal/repos"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret string
	repo      *repos.ConnRepo
}

func NewAuthService(conn *repos.ConnRepo, jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
		repo:      conn,
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

// Функция регистрации пользователя
func (auth *AuthService) Register(login, password string, ctx context.Context) (string, error) {
	id, err := auth.repo.Register(login, password, ctx)
	if err != nil {
		return "", err
	}
	token, err := auth.CreateToken(id)
	if err != nil {
		return "", err
	}
	return token, nil
}
