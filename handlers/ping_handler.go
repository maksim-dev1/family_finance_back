package handlers

import "net/http"

// PingHandler отвечает на запросы /ping, возвращая "pong".
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
