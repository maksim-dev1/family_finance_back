package handlers

import (
	"encoding/json"
	"net/http"

	"family_finance_back/internal/service"
)

// UserHandler обрабатывает HTTP запросы, связанные с данными пользователя
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler создает новый экземпляр UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// UpdateUserRequest представляет запрос на обновление данных пользователя
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Nickname string `json:"nickname"`
}

// GetUserHandler обрабатывает запрос на получение данных пользователя
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка Authorization
	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", "Токен не предоставлен")
		return
	}

	// Убираем префикс "Bearer " если он есть
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Получаем данные пользователя
	user, err := h.userService.GetUserByToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", err.Error())
		return
	}

	// Отправляем данные пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUserHandler обрабатывает запрос на обновление данных пользователя
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из заголовка Authorization
	token := r.Header.Get("Authorization")
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", "Токен не предоставлен")
		return
	}

	// Убираем префикс "Bearer " если он есть
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	// Получаем текущего пользователя
	user, err := h.userService.GetUserByToken(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Ошибка авторизации", err.Error())
		return
	}

	// Парсим данные для обновления
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Не удалось прочитать данные запроса")
		return
	}

	// Обновляем только предоставленные поля
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Surname != "" {
		user.Surname = req.Surname
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	// Сохраняем изменения
	if err := h.userService.UpdateUser(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка обновления", err.Error())
		return
	}

	// Отправляем обновленные данные пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// SearchUserByEmailHandler обрабатывает запрос на поиск пользователя по email
func (h *UserHandler) SearchUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем email из query параметра
	email := r.URL.Query().Get("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "Неверный запрос", "Параметр 'email' обязателен")
		return
	}

	// Получаем данные пользователя
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден", err.Error())
		return
	}

	// Отправляем данные пользователя
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
