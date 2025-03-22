package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	uuid, err := uuid.Parse(idString)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Failed to convert UUID string too UUID"})
		return
	}

	dbChirp, err := cfg.dbQueries.GetOneChirp(r.Context(), uuid)
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, errResponse{Error: "Data base request for Chirp failed"})
		return
	}

	response := ChirpResponse{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbChirp.UpdatedAt.Format(time.RFC3339),
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}

	writeJSONResponse(w, http.StatusOK, response)
}