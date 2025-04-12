package handlers

import (
	"encoding/json"
	"net/http"

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

// RequestRegistrationCodeHandler для запроса кода регистрации
func (h *AuthHandler) RequestRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Name == "" || req.Surname == "" {
		http.Error(w, "неправильный запрос", http.StatusBadRequest)
		return
	}
	err := h.authService.RequestRegistrationCode(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("код отправлен"))
}

// VerifyRegistrationCodeHandler для подтверждения регистрации
func (h *AuthHandler) VerifyRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Code == "" || req.Name == "" || req.Surname == "" {
		http.Error(w, "неправильный запрос", http.StatusBadRequest)
		return
	}
	err := h.authService.VerifyRegistrationCode(req.Email, req.Code, req.Name, req.Surname, req.Nickname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("пользователь создан"))
}
