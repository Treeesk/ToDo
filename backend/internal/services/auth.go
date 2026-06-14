package services

// Создание и проверка JWT токенов. Регистрация и логин пользователей.

import (
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

// Функция по созданию JWT
func (auth *AuthService) CreateToken(user_id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": user_id,
			"iat": jwt.NewNumericDate(time.Now()),
			"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		})
	tokenString, err := token.SignedString(auth.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Функция проверки JWT
func (auth *AuthService) VerifyToken(tokenString string) error {

}
