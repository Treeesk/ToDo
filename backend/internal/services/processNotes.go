package services

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
