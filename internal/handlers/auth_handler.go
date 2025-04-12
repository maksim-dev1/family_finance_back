package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"family_finance_back/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest принимает JSON с полем email для запроса кода
type LoginRequest struct {
	Email string `json:"email"`
}

// CodeRequest принимает JSON с email и кодом
type CodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

// RegistrationRequest для начала регистрации
type RegistrationRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

// RegisterCodeRequest для подтверждения регистрации
type RegisterCodeRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

// RequestLoginCodeHandler для запроса кода входа
func (h *AuthHandler) RequestLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "неправильный запрос", http.StatusBadRequest)
		return
	}
	err := h.authService.RequestLoginCode(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("код отправлен"))
}

// VerifyLoginCodeHandler для проверки кода и авторизации
func (h *AuthHandler) VerifyLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Code == "" {
		http.Error(w, "неправильный запрос", http.StatusBadRequest)
		return
	}
	token, err := h.authService.VerifyLoginCode(req.Email, req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Возвращаем токен
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RequestRegistrationCodeHandler для запроса кода регистрации с возвратом UUID
func (h *AuthHandler) RequestRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
    type Request struct {
        Email string `json:"email"`
    }
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
        http.Error(w, "неправильный запрос", http.StatusBadRequest)
        return
    }
    // Получаем UUID от сервиса
    uuidKey, err := h.authService.RequestRegistrationCode(req.Email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    // Отправляем UUID клиенту
    json.NewEncoder(w).Encode(map[string]string{"registration_id": uuidKey})
}


// VerifyRegistrationCodeHandler для подтверждения регистрации с использованием UUID
func (h *AuthHandler) VerifyRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
    type Request struct {
        RegistrationID string `json:"registration_id"`
        Code           string `json:"code"`
        Name           string `json:"name"`
        Surname        string `json:"surname"`
        Nickname       string `json:"nickname"`
    }
    var req Request
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RegistrationID == "" || req.Code == "" || req.Name == "" || req.Surname == "" {
        http.Error(w, "неправильный запрос", http.StatusBadRequest)
        return
    }

    err := h.authService.VerifyRegistrationCode(req.RegistrationID, req.Code, req.Name, req.Surname, req.Nickname)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("пользователь создан"))
}


func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем токен из заголовка Authorization в формате "Bearer <token>"
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Неверный формат заголовка Authorization", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	err := h.authService.Logout(token)
	if err != nil {
		http.Error(w, "Ошибка логаута: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Можно вернуть результат логаута
	json.NewEncoder(w).Encode(map[string]string{"message": "Вы успешно вышли из аккаунта"})
}
