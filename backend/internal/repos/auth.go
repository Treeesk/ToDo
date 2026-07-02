package repos

// Методы для взаимодействия с бд при аутентификации
import (
	"ProjectGo/backend/internal/customerrors"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Создание пользователя, возвращаем id и рефреш токен в случае успеха
func (repo *ConnRepo) Register(login, password string, ctx context.Context, exp_refresh time.Time) (int, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, "", err // слишком длинный или короткий пароль
	}
	var userId int
	tx, err := repo.Conn.Begin(ctx) // начинаем транзакцию
	if err != nil {
		return -1, "", err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, "INSERT INTO users (user_login, user_password) VALUES ($1, $2) RETURNING id", login, string(hash)).Scan(&userId)
	if err != nil {
		return -1, "", err
	}
	refresh_token := create_refresh_token()
	hash_refresh_token := sha256.Sum256([]byte(refresh_token))
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES($1, $2, $3, $4)", userId, hash_refresh_token[:], exp_refresh, time.Now())
	if err != nil {
		return -1, "", err
	}
	err = tx.Commit(ctx) // заканчиваем транзакцию
	if err != nil {
		return -1, "", err
	}
	return userId, refresh_token, nil
}

// Логин пользователя, возвращаем id и рефреш токен в случае успеха
func (repo *ConnRepo) Login(login, password string, ctx context.Context, exp_refresh time.Time) (int, string, error) {
	var userId int
	var hashpass []byte
	tx, err := repo.Conn.Begin(ctx) // начинаем транзакцию
	if err != nil {
		return -1, "", err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, "SELECT id, user_password FROM users WHERE user_login = $1", login).Scan(&userId, &hashpass)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, "", &customerrors.UserError{What: fmt.Sprintf("invalid login or password, login: %s", login)}
		}
		return -1, "", err
	}
	err = bcrypt.CompareHashAndPassword(hashpass, []byte(password))
	if err != nil {
		return -1, "", &customerrors.UserError{What: fmt.Sprintf("invalid login or password, login: %s", login)}
	}
	// Удаляем старые токены (если их больше 100 на одного юзера)
	_, err = tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE id IN (SELECT id FROM refresh_tokens WHERE user_id = $1 ORDER BY created_at ASC LIMIT(SELECT GREATEST(COUNT(*) - 99, 0) FROM refresh_tokens WHERE user_id = $1))", userId)
	if err != nil {
		return -1, "", err
	}
	refresh_token := create_refresh_token()
	hash_refresh_token := sha256.Sum256([]byte(refresh_token))
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES($1, $2, $3, $4)", userId, hash_refresh_token[:], exp_refresh, time.Now())
	if err != nil {
		return -1, "", err
	}
	err = tx.Commit(ctx) // заканчиваем транзакцию
	if err != nil {
		return -1, "", err
	}
	return userId, refresh_token, nil
}

// Функция удаления refresh_token из DB
func (repo *ConnRepo) LogOut(ctx context.Context, refresh_token string) error {
	token_hash := sha256.Sum256([]byte(refresh_token))
	_, err := repo.Conn.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1", token_hash[:])
	return err
}

// Функция проверки refresh token пользователя и в случае успеха создания нового
func (repo *ConnRepo) Refresh(refresh string, exp_refresh time.Time, ctx context.Context) (int, string, error) {
	type Data struct {
		UserId     int
		Expires_at time.Time // чтобы проверить не просрочился ли refresh токен
	}
	var data Data
	token_hash := sha256.Sum256([]byte(refresh))
	tx, err := repo.Conn.Begin(ctx)
	if err != nil {
		return -1, "", err
	}
	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, "SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1", token_hash[:]).Scan(&data.UserId, &data.Expires_at)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, "", &customerrors.UserError{What: "unknown refresh-token"}
		}
		return -1, "", err
	}
	tag, err := tx.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash = $1", token_hash[:]) // удаляем старый refresh токен
	if err != nil {
		return -1, "", err
	}
	// если пользователь конкурентно уже удалил токен, мы не должны еще один создавать
	if tag.RowsAffected() == 0 {
		return -1, "", &customerrors.UserError{What: "unknown refresh-token"}
	}
	// токен просрочен
	if time.Now().After(data.Expires_at) {
		return -1, "", &customerrors.UserError{What: "refresh token expired"}
	}
	// создаем новый refresh токен и помещаем его в бд
	refresh_token := create_refresh_token()
	hash_refresh_token := sha256.Sum256([]byte(refresh_token))
	_, err = tx.Exec(ctx, "INSERT INTO refresh_tokens (user_id, token_hash, expires_at, created_at) VALUES($1, $2, $3, $4)", data.UserId, hash_refresh_token[:], exp_refresh, time.Now())
	if err != nil {
		return -1, "", err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return -1, "", err
	}
	return data.UserId, refresh_token, nil
}

// Создание refresh токена. Возвращаем строку нехешированную
func create_refresh_token() string {
	b := make([]byte, 32)
	rand.Read(b)
	refresh_token := base64.RawURLEncoding.EncodeToString(b)
	return refresh_token
}
