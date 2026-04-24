package services

import (
	"ProjectGo/backend/internal/entity"
	"ProjectGo/backend/internal/repos"
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
func (s *NotesStore) GetAll(user_id int) ([]entity.Note, error) {
	return s.repo.GetAllNotes(user_id)
}

// Добавляем заметку
func (s *NotesStore) Add(user_id int, text string) error {
	return s.repo.AddNotedb(user_id, text)
}

// Удаление заметки
func (s *NotesStore) Del(user_id, id int) error {
	if id < 0 {
		return fmt.Errorf("bad ID: %d", id)
	}
	if user_id < 0 {
		return fmt.Errorf("bad ID: %d", user_id)
	}
	return s.repo.DeleteNotedb(user_id, id)
}

// Изменение заметки
func (s *NotesStore) Edit(user_id, id int, text string) error {
	if id < 0 {
		return fmt.Errorf("bad ID: %d", id)
	}
	if user_id < 0 {
		return fmt.Errorf("bad ID: %d", user_id)
	}
	return s.repo.EditNotedb(user_id, id, text)
}
