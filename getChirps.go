package main

import (
	"net/http"
	"time"
)


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Failed to fetch chirps"})
	}

	responseList := []ChirpResponse{}

	for _, dbChirp := range chirps {
		responseList = append(responseList, ChirpResponse{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt.Format(time.RFC3339),
			UpdatedAt: dbChirp.UpdatedAt.Format(time.RFC3339),
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	writeJSONResponse(w, http.StatusOK, responseList)
 }