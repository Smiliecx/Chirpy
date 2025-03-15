package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Smiliecx/Chirpy/internal/database"
	"github.com/google/uuid"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		 Body string `json:"body"`
		 User_id string `json:"user_id"`
	}
 
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()
 
	decoder := json.NewDecoder(r.Body)
	var newChirp chirp
	err := decoder.Decode(&newChirp)
 
	if err != nil {
		 writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Something went wrong"})
		 return
	}
 
	if len(newChirp.Body) > 140 {
		 writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Chirp is too long"})
		 return
	}
 
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Fields(newChirp.Body)
	var cleanWords []string
 
	for _, word := range words {
		 if contains(profaneWords, strings.ToLower(word)) {
			 cleanWords = append(cleanWords, "****")
		 } else {
			 cleanWords = append(cleanWords, word)
		 }
	}
 
	cleanBody := strings.Join(cleanWords, " ")
	cleanBody = strings.TrimSpace(cleanBody)

	parsedUUID, err := uuid.Parse(newChirp.User_id)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Invalid UUID For User"})
		return
	}

	dbChirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanBody, UserID: parsedUUID})
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Chirp failed insertion into DB"})
		return
	}

	response := ChirpResponse{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbChirp.UpdatedAt.Format(time.RFC3339),
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
 
	writeJSONResponse(w, http.StatusCreated, response)
 }
 
 func contains(list []string, item string) bool {
	 for _, v := range list {
		 if v == item {
			 return true
		 }
	 }
	 return false
 }