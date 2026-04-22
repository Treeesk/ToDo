package repos

// Для взаимодействия с Базой данных

import (
	"ProjectGo/backend/internal/config"
	"ProjectGo/backend/internal/entity"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnRepo struct {
	Conn *pgxpool.Pool
}

// Подключение к базе данных
func ConnUrlRepos(cfg *config.Config) *ConnRepo {
	pool, err := pgxpool.New(context.TODO(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	return &ConnRepo{Conn: pool}
}

// Получение всех заметок из базы данных(пока что всех пользователей)
func (repo *ConnRepo) GetAllNotes() ([]entity.Note, error) {
	rows, err := repo.Conn.Query(context.TODO(), "SELECT * FROM notes")
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return []entity.Note{}, fmt.Errorf("error on select rows: %d", err)
	}
	var notes []entity.Note
	for rows.Next() {
		var note entity.Note
		rows.Scan(&note.ID, &note.User_id, &note.Text)
		notes = append(notes, note)
	}
	defer rows.Close()
	return notes, nil
}

// Добавление заметки в бд
func (repo *ConnRepo) AddNotebd(user_id int, text string) error {
	_, err := repo.Conn.Exec(context.TODO(), "INSERT INTO notes (user_id, note) VALUES($1, $2)", user_id, text)
	if err != nil {
		return err
	}
	return nil
}
