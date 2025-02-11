package handlers

import (
	"net/http"

	"family_finance_back/repository"
)

// UserHandler содержит ссылку на репозиторий пользователей.
type UserHandler struct {
	userRepo repository.UserRepository
}

// NewUserHandler создаёт новый экземпляр UserHandler.
func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

// GetUser обрабатывает GET-запрос для получения данных текущего пользователя.
// Для получения email используется контекст, куда его ранее записал middleware TokenAuthMiddleware.
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Извлекаем email из контекста.
	email, ok := r.Context().Value(ContextKeyUserEmail).(string)
	if !ok || email == "" {
		respondWithError(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	// Получаем данные пользователя по email.
	user, err := h.userRepo.GetUserByEmail(email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка при получении пользователя")
		return
	}
	if user == nil {
		respondWithError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}

	// Возвращаем данные пользователя в формате JSON.
	respondWithJSON(w, http.StatusOK, user)
}
