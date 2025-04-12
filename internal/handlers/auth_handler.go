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

// Для запроса кода авторизации (login) – требуется только email
type LoginRequest struct {
	Email string `json:"email"`
}

// Для подтверждения кода авторизации – теперь передаётся временный UUID (temp_id)
type VerifyLoginRequest struct {
	TempID string `json:"temp_id"`
	Code   string `json:"code"`
}

// Для запроса кода регистрации – требуется только email
type RegistrationRequest struct {
	Email string `json:"email"`
}

// Для подтверждения регистрации – передаются temp_id, код, имя, фамилия и никнейм
type VerifyRegistrationRequest struct {
	TempID   string `json:"temp_id"`
	Code     string `json:"code"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

// RequestLoginCodeHandler для запроса кода входа
func (h *AuthHandler) RequestLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Неправильный запрос: поле 'email' обязательно для заполнения", http.StatusBadRequest)
		return
	}
	tempID, err := h.authService.RequestLoginCode(req.Email)
	if err != nil {
		http.Error(w, "Ошибка запроса кода авторизации: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"temp_id": tempID})
}

// VerifyLoginCodeHandler для проверки кода и авторизации
func (h *AuthHandler) VerifyLoginCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TempID == "" || req.Code == "" {
		http.Error(w, "Неправильный запрос: поля 'temp_id' и 'code' обязательны для заполнения", http.StatusBadRequest)
		return
	}
	token, err := h.authService.VerifyLoginCode(req.TempID, req.Code)
	if err != nil {
		http.Error(w, "Ошибка верификации кода авторизации: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RequestRegistrationCodeHandler изменён: принимает только email и возвращает temp_id
func (h *AuthHandler) RequestRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Неправильный запрос: поле 'email' обязательно для заполнения", http.StatusBadRequest)
		return
	}

	// Метод возвращает temp_id или подробную ошибку
	tempID, err := h.authService.RequestRegistrationCode(req.Email)
	if err != nil {
		// Клиенту возвращаем детальное сообщение об ошибке
		http.Error(w, "Ошибка запроса кода регистрации: "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"temp_id": tempID})
}

// VerifyRegistrationCodeHandler изменён: принимает temp_id и дополнительные данные
func (h *AuthHandler) VerifyRegistrationCodeHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil ||
		req.TempID == "" || req.Code == "" || req.Name == "" || req.Surname == "" {
		http.Error(w, "Неправильный запрос: поля 'temp_id', 'code', 'name' и 'surname' обязательны", http.StatusBadRequest)
		return
	}
	err := h.authService.VerifyRegistrationCode(req.TempID, req.Code, req.Name, req.Surname, req.Nickname)
	if err != nil {
		// Более подробное сообщение об ошибке пользователю
		http.Error(w, "Ошибка подтверждения регистрации: "+err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Пользователь успешно создан"))
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
