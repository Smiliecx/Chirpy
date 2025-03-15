package main

import (
	"encoding/json"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}