package services

import (
	"ProjectGo/backend/internal/customerrors"
	"ProjectGo/backend/internal/entity"
	"ProjectGo/backend/internal/repos"
	"context"
	"fmt"
)

type NotesStore struct {
	repo *repos.ConnRepo
}

// Сохранения подключения к бд, чтобы можно было вызывать к этому подключению различные методы
func NewNotesStore(conn *repos.ConnRepo) *NotesStore {
	return &NotesStore{
		repo: conn,
	}
}

// Возвращаем все заметки
func (s *NotesStore) GetAll(ctx context.Context, user_id int) ([]entity.Note, error) {
	return s.repo.GetAllNotes(ctx, user_id)
}

// Добавляем заметку
func (s *NotesStore) Add(ctx context.Context, user_id int, text string) error {
	return s.repo.AddNotedb(ctx, user_id, text)
}

// Удаление заметки
func (s *NotesStore) Del(ctx context.Context, user_id, id int) error {
	if id < 0 {
		return fmt.Errorf("bad ID: %d", id)
	}
	if user_id < 0 {
		return fmt.Errorf("bad ID: %d", user_id)
	}
	return s.repo.DeleteNotedb(ctx, user_id, id)
}

// Изменение заметки
func (s *NotesStore) Edit(ctx context.Context, user_id, id int, text string) error {
	if id < 0 {
		return &customerrors.ErrorNotFound{What: "note not found", Id: id, User_id: user_id}
	}
	if user_id < 0 {
		return &customerrors.ErrorNotFound{What: "note not found", Id: id, User_id: user_id}
	}
	return s.repo.EditNotedb(ctx, user_id, id, text)
}
