package handlers

// Хэндлеры для работы с пользователем(login, register, logout)
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Хэндлер для регистрации пользователя
func (h *HandlerNotes) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
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
	if us.Login == "" || us.Password == "" {
		writeJsonError(w, http.StatusBadRequest, "invalid login or password")
		log.Printf("invalid login: %s or password: %s", us.Login, us.Password)
		return
	}
	expires_access := time.Now().Add(time.Minute * 15) // время жизни куки
	expires_refresh := time.Now().AddDate(0, 0, 30)    // время жизни refresh токена
	access_token, refresh_token, err := h.authService.Register(us.Login, us.Password, ctx, expires_access, expires_refresh)
	if err != nil {
		HandleError(w, err)
		log.Println("Register failed:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    access_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_access,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refresh_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_refresh,
	})
	w.WriteHeader(http.StatusCreated)
}

// Хэндлер для логина пользователя
func (h *HandlerNotes) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
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
	expires_access := time.Now().Add(time.Minute * 15) // время жизни куки
	expires_refresh := time.Now().AddDate(0, 0, 30)    // время жизни refresh токена
	access_token, refresh_token, err := h.authService.Login(us.Login, us.Password, ctx, expires_access, expires_refresh)
	if err != nil {
		HandleError(w, err)
		log.Println("login failed:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    access_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_access,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refresh_token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_refresh,
	})
	w.WriteHeader(http.StatusOK)
}

// Функция выхода пользователя из своего профиля
func (h *HandlerNotes) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // работаем с контекстом(пользователь может закрыть соединение или мы будем долго выполнять работу)
	defer cancel()
	cook, err := r.Cookie("refresh-token")
	if err != nil {
		writeJsonError(w, http.StatusUnauthorized, "Unauthorized")
		log.Println("Unauthorized user")
		return
	}
	err = h.authService.LogOut(ctx, cook.Value)
	if err != nil {
		HandleError(w, err)
		log.Println("login failed:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}

// Хэндлер для обновления access и refresh токенов
func (h *HandlerNotes) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonError(w, http.StatusMethodNotAllowed, "Method not allowed")
		log.Println("error: method not allowed")
		return
	}
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
		HandleError(w, err)
		// удаление старых кук
		http.SetCookie(w, &http.Cookie{
			Name:     "access-token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh-token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		log.Println("refresh failed:", err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    access,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_access,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh-token",
		Value:    refresh,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires_refresh,
	})
	w.WriteHeader(http.StatusOK)
}
