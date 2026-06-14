package repos

// Для взаимодействия с Базой данных

import (
	"ProjectGo/backend/internal/config"
	"ProjectGo/backend/internal/customerrors"
	"ProjectGo/backend/internal/entity"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnRepo struct {
	Conn *pgxpool.Pool
}

// Подключение к базе данных
func ConnUrlRepos(ctx context.Context, cfg *config.Config) *ConnRepo {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return &ConnRepo{Conn: pool}
}

// Получение всех заметок из базы данных
func (repo *ConnRepo) GetAllNotes(ctx context.Context, user_id int) ([]entity.Note, error) {
	rows, err := repo.Conn.Query(ctx, "SELECT * FROM notes WHERE user_id = $1", user_id)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()
	var notes []entity.Note
	for rows.Next() {
		var note entity.Note
		if err := rows.Scan(&note.ID, &note.User_id, &note.Text); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if notes == nil {
		return []entity.Note{}, nil
	}
	return notes, nil
}

// Добавление заметки в бд
func (repo *ConnRepo) AddNotedb(ctx context.Context, user_id int, text string) error {
	_, err := repo.Conn.Exec(ctx, "INSERT INTO notes (user_id, note) VALUES($1, $2)", user_id, text)
	if err != nil {
		return err
	}
	return nil
}

// Удаление заметки из бд
func (repo *ConnRepo) DeleteNotedb(ctx context.Context, user_id, id int) error {
	tag, err := repo.Conn.Exec(ctx, "DELETE FROM notes WHERE user_id = $1 AND id = $2", user_id, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return &customerrors.ErrorNotFound{What: "note not found", Id: id, User_id: user_id}
	}
	return nil
}

// Редактирование заметки в бд
func (repo *ConnRepo) EditNotedb(ctx context.Context, user_id, id int, text string) error {
	tag, err := repo.Conn.Exec(ctx, "UPDATE notes SET note = $1 WHERE user_id = $2 AND id = $3", text, user_id, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return &customerrors.ErrorNotFound{What: "note not found", Id: id, User_id: user_id}
	}
	return nil
}
