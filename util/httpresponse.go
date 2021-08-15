package util

import (
	"encoding/json"
	"net/http"
)

func RespondJson(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
