package customerrors

import "fmt"

type ErrorNotFound struct {
	What    string
	Id      int
	User_id int
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("%s, id: %d, user_id: %d", e.What, e.Id, e.User_id)
}
