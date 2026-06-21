package repos

// Методы для взаимодействия с бд при аунтентификации
import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Создание пользователя, возвращаем id
func (repo *ConnRepo) Register(login, password string, ctx context.Context) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("error: %v", err) // слишком длинный или короткий пароль
	}
	var userId int
	err = repo.Conn.QueryRow(ctx, "INSERT INTO users (user_login, user_password) VALUES ($1, $2) RETURNING id", login, string(hash)).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("error: %v", err)
	}
	return userId, nil
}
