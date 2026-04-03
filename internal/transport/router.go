package transport

import (
	"ProjectGo/internal/handlers"
	"net/http"
)

func Setuprouter() {
	http.HandleFunc("/", handlers.HandlerStart)
}
