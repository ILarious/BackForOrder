package handler

import (
	"encoding/json"
	"net/http"
)

type Health struct {
	Status string `json:"status"`
}

func NewHealth() *Health {
	return &Health{}
}

func (h Health) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(Health{Status: "OK"})
}
