package services

type Note struct {
	id   int
	text string
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

func (s *NotesStore) GetAll() []Note {
	return s.notes
}
