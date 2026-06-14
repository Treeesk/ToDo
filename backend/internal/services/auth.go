package services

// Создание и проверка JWT токенов

type AuthService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: jwtSecret,
	}
}

// Функция по созданию JWT
func (auth *AuthService) CreateToken() (string, error) {

}
