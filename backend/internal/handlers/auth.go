package handlers

// Хэндлеры для работы с пользователем(login, register, logout)
import (
	"ProjectGo/backend/internal/customerrors"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Хэндлер для регистрации пользователя
func (h *HandlerNotes) Register(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()
	var us user
	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	expires := time.Now().Add(time.Minute * 15) // время жизни куки
	token, err := h.authService.Register(us.Login, us.Password, ctx, expires)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}
		if errors.Is(err, customerrors.ErrTokenCreate) {
			writeJsonError(w, http.StatusInternalServerError, "server error")
			log.Println(err)
			return
		}
		if errors.Is(err, bcrypt.ErrHashTooShort) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
			writeJsonError(w, http.StatusBadRequest, "too long or too short password")
			log.Println(err)
			return
		}
		HandleError(w, err) // бизнес и бд ошибки
		log.Println("database error:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Хэндлер для логина пользователя
func (h *HandlerNotes) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()
	var us user
	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		jsonDecodeError(w, err)
		log.Println("decode error: ", err)
		return
	}
	expires := time.Now().Add(time.Minute * 15) // время жизни куки
	token, err := h.authService.Login(us.Login, us.Password, ctx, expires)
	if err != nil {
		// Серверный ошибки
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			writeJsonError(w, http.StatusGatewayTimeout, "request timeout")
			log.Println(err)
			return
		}
		if errors.Is(err, customerrors.ErrTokenCreate) {
			writeJsonError(w, http.StatusInternalServerError, "server error")
			log.Println(err)
			return
		}
		HandleError(w, err) // бизнес и бд ошибки
		log.Println("database error:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Функция выхода пользователя из своего профиля
func (h *HandlerNotes) LogOut(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()

}

// Хэндлер для обновления access и refresh токенов
func (h *HandlerNotes) Refresh(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()
	cook, err := r.Cookie("refresh-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized")
		log.Println("Unauthorized user")
		return
	}
	expires_access := time.Now().Add(time.Minute * 15) // время жизни куки
	expires_refresh := time.Now().AddDate(0, 0, 30)    // время жизни refresh токена
	access, refresh, err := h.authService.Refresh(cook.Value, ctx, expires_access, expires_refresh)
	if err != nil {

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
