// handlers/response.go
package handlers

import (
	"encoding/json"
	"net/http"
)

// respondWithError отправляет JSON-ответ с сообщением об ошибке и нужным HTTP-кодом.
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	// Можно возвращать объект с полем error и, например, timestamp, если нужно.
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

// respondWithJSON отправляет успешный ответ с данными в формате JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
