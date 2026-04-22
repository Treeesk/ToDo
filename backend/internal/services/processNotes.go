package services

import (
	"ProjectGo/backend/internal/entity"
	"ProjectGo/backend/internal/repos"
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
	return s.repo.AddNotebd(user_id, text)
}

// Удаление заметки
func (s *NotesStore) Del(id int) error {
	// if id < 1 || id >= s.nextID {
	// 	return fmt.Errorf("bad ID: %d", id)
	// }
	// s.notes = slices.Delete(s.notes, id-1, id)
	// for i := id - 1; i < len(s.notes); i++ {
	// 	s.notes[i].ID--
	// }
	return nil
}

// Изменение заметки
func (s *NotesStore) Edit(id int, text string) error {
	// if id < 1 || id >= s.nextID {
	// 	return fmt.Errorf("bad ID: %d", id)
	// }
	// s.notes[id-1].Text = text
	return nil
}
