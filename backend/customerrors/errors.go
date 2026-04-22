package customerrors

import "fmt"

type ErrorNotFound struct {
	What string
	Id   int
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("%s, id: %d", e.What, e.Id)
}
