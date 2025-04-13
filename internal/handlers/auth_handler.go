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

// Вспомогательная функция для возврата ошибок в формате JSON
func respondWithError(w http.ResponseWriter, statusCode int, errMsg, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errMsg,
		"details": details,
	})
}

type LoginRequest struct {
	Email string `json:"email"`
}

type VerifyLoginRequest struct {
	TempID string `json:"temp_id"`
	Code   string `json:"code"`
}

type RegistrationRequest struct {
	Email string `json:"email"`
}

type VerifyRegistrationRequest struct {
	TempID   string `json:"temp_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

func (h *AuthHandler) RequestLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Поле 'email' обязательно для заполнения")
		return
	}
	tempID, err := h.authService.RequestLoginCode(req.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Ошибка запроса кода авторизации", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"temp_id": tempID})
}

func (h *AuthHandler) VerifyLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TempID == "" || req.Code == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Поля 'temp_id' и 'code' обязательны для заполнения")
		return
	}
	token, err := h.authService.VerifyLoginCode(req.TempID, req.Code)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Ошибка верификации кода авторизации", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) RequestRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Поле 'email' обязательно для заполнения")
		return
	}

	tempID, err := h.authService.RequestRegistrationCode(req.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Ошибка запроса кода регистрации", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"temp_id": tempID})
}

func (h *AuthHandler) VerifyRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		req.TempID == "" || req.Code == "" || req.Name == "" || req.Surname == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Поля 'temp_id', 'code', 'name' и 'surname' обязательны")
		return
	}
	// Получаем токен после успешной регистрации
	token, err := h.authService.VerifyRegistrationCode(req.TempID, req.Code, req.Name, req.Surname, req.Nickname)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Ошибка подтверждения регистрации", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", "Отсутствует заголовок Authorization")
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", "Неверный формат заголовка Authorization")
		return
	}

	token := parts[1]
	err := h.authService.Logout(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка логаута", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Вы успешно вышли из аккаунта"})
}
