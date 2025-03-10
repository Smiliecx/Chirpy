package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirp(w http.ResponseWriter, r *http.Request) {
   type chirp struct {
		Body string `json:"body"`
   }

   type errResponse struct {
		Error string `json:"error"`
   }

   type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
   }

   w.Header().Set("Content-Type", "application/json")
   defer r.Body.Close()

   decoder := json.NewDecoder(r.Body)
   var newChirp chirp
   err := decoder.Decode(&newChirp)

   if err != nil {
		resp := errResponse{Error: "Something went wrong"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
   }

   if len(newChirp.Body) > 140 {
		resp := errResponse{Error: "Chirp is too long"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
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

   resp := validResponse{CleanedBody: cleanBody}
   w.WriteHeader(http.StatusOK)
   json.NewEncoder(w).Encode(resp)
}

func contains(list []string, item string) bool {
    for _, v := range list {
        if v == item {
            return true
        }
    }
    return false
}