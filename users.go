package main

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Smiliecx/Chirpy/internal/auth"
	"github.com/Smiliecx/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type user_request struct {
		Email string `json:"email"`
		Password string `json: "password"`
	}

	type user_response struct {
		Id string `json:"id"`
		Created_at string `json:"created_at"`
		Updated_at string `json:"updated_at"`
		Email string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
   defer r.Body.Close()

   decoder := json.NewDecoder(r.Body)
   var user user_request
   err := decoder.Decode(&user)
   if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Something went wrong while Decoding user"})
		return
   }

   if !validateEmail(user.Email) {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Email field was not valid"})
		return
   }

   hashed_password, err := auth.HashPassword(user.Password)
   if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Failed to hash password for user"})
		return
   }

   new_user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hashed_password,
		Email: user.Email,
   })
   
   if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "User Query Failed"})
		return
   }

   

   writeJSONResponse(w, http.StatusCreated, user_response{
		Id: new_user.ID.String(),
		Created_at: new_user.CreatedAt.String(),
		Updated_at: new_user.UpdatedAt.String(),
		Email: new_user.Email,
   })
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[^@]+@[^@]+\.[^@]+$`)
	if !re.MatchString(email) {
		return false
	}
	return true
}