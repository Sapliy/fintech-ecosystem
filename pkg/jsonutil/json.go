package jsonutil

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorJSON(w http.ResponseWriter, errMsg string) {
	log.Fatalf("Error: %s", errMsg)
	WriteJSON(w, http.StatusBadRequest, map[string]string{"error": errMsg})
}
