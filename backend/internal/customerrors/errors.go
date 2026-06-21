package customerrors

import (
	"errors"
	"fmt"
)

// JWT обертки ошибок
var (
	ErrTokenCreate = errors.New("token creation failed")
	ErrTokenParse  = errors.New("token parse failed")
)

type ErrorNotFound struct {
	What    string
	Id      int
	User_id int
}

type UserError struct {
	What  string
	Login string
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("%s, id: %d, user_id: %d", e.What, e.Id, e.User_id)
}

func (e *UserError) Error() string {
	return fmt.Sprintf("error: %s, login: %s", e.What, e.Login)
}
