package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"family_finance_back/internal/service"
)

// AuthHandler обрабатывает HTTP запросы, связанные с авторизацией
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// respondWithError отправляет ответ с ошибкой в формате JSON
func respondWithError(w http.ResponseWriter, statusCode int, errMsg, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errMsg,
		"details": details,
	})
}

// LoginRequest представляет запрос на получение кода для входа
type LoginRequest struct {
	Email string `json:"email"`
}

// VerifyLoginRequest представляет запрос на проверку кода входа
type VerifyLoginRequest struct {
	TempID string `json:"temp_id"`
	Code   string `json:"code"`
}

// RegistrationRequest представляет запрос на получение кода для регистрации
type RegistrationRequest struct {
	Email string `json:"email"`
}

// VerifyRegistrationRequest представляет запрос на проверку кода регистрации
type VerifyRegistrationRequest struct {
	TempID   string `json:"temp_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

// RequestLoginCodeHandler обрабатывает запрос на получение кода для входа
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

// VerifyLoginCodeHandler обрабатывает запрос на проверку кода входа
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

// RequestRegistrationCodeHandler обрабатывает запрос на получение кода для регистрации
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

// VerifyRegistrationCodeHandler обрабатывает запрос на проверку кода регистрации
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

// LogoutHandler обрабатывает запрос на выход из системы
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
