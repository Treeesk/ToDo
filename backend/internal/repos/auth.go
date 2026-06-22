package repos

// Методы для взаимодействия с бд при аунтентификации
import (
	"ProjectGo/backend/internal/customerrors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Создание пользователя, возвращаем id в случае успеха
func (repo *ConnRepo) Register(login, password string, ctx context.Context) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, err // слишком длинный или короткий пароль
	}
	var userId int
	err = repo.Conn.QueryRow(ctx, "INSERT INTO users (user_login, user_password) VALUES ($1, $2) RETURNING id", login, string(hash)).Scan(&userId)
	if err != nil {
		return -1, err
	}
	return userId, nil
}

// Логин пользователя, возвращаем id в случае успеха
func (repo *ConnRepo) Login(login, password string, ctx context.Context) (int, error) {
	var userId int
	var hashpass []byte
	err := repo.Conn.QueryRow(ctx, "SELECT id, user_password FROM users WHERE user_login = $1", login).Scan(&userId, &hashpass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, &customerrors.UserError{What: "invalid login or password", Login: login}
		}
		return -1, err
	}
	err = bcrypt.CompareHashAndPassword(hashpass, []byte(password))
	if err != nil {
		return -1, &customerrors.UserError{What: "invalid login or password", Login: login}
	}
	return userId, nil
}
