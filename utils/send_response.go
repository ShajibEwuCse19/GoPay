package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendError writes an error in one consistent JSON shape: {"error": "message"}.
func SendError(w http.ResponseWriter, status int, message string) {
	SendResponse(w, status, map[string]string{"error": message})
}
