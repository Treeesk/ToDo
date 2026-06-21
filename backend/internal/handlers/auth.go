package handlers

// Хэндлеры для работы с пользователем(login, register, logout)
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// Хэндлер для регистрации пользователя
func (h *HandlerNotes) Register(w http.ResponseWriter, r *http.Request) {
	type user struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()
	var us user
	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	token, err := h.authService.Register(us.Login, us.Password, ctx)

}
