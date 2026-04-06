package services

import (
	"fmt"
	"slices"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type NotesStore struct {
	nextID int
	notes  []Note
}

func NewNotesStore() *NotesStore {
	return &NotesStore{
		nextID: 1,
		notes:  []Note{},
	}
}

// Возвращаем все заметки
func (s *NotesStore) GetAll() []Note {
	return s.notes
}

// Добавляем заметку
func (s *NotesStore) Add(text string) {
	note := Note{ID: s.nextID, Text: text}
	s.nextID++
	s.notes = append(s.notes, note)
}

// Удаление заметки
func (s *NotesStore) Del(id int) error {
	if id < 1 || id >= s.nextID {
		return fmt.Errorf("bad ID: %d", id)
	}
	s.notes = slices.Delete(s.notes, id-1, id)
	for i := id - 1; i < len(s.notes); i++ {
		s.notes[i].ID--
	}
	return nil
}

// Изменение заметки
func (s *NotesStore) Edit(id int, text string) error {
	if id < 1 || id >= s.nextID {
		return fmt.Errorf("bad ID: %d", id)
	}
	s.notes[id-1].Text = text
	return nil
}
