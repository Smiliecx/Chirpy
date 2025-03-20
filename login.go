package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Smiliecx/Chirpy/internal/auth"
	"github.com/Smiliecx/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
   	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	var result loginRequest
	err := decoder.Decode(&result)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Something went wrong while Decoding user"})
		return
	}

	new_user, err := cfg.dbQueries.FindUserByEmail(r.Context(), result.Email)
	if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Could not find user with Email"})
			return
	}

	validPassword := auth.CheckPasswordHash(result.Password, new_user.HashedPassword)
	if validPassword != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Password is invalid"})
		return
	}

	jwt, err := auth.MakeJWT(new_user.ID, cfg.authSecret, time.Duration(time.Hour))
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Could not generate jwt"})
		return
   }

   refresh_token, _ := auth.MakeRefreshToken()
   _, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refresh_token,
		UserID: new_user.ID,
		ExpiresAt: time.Now().Add(60 * time.Hour * 24),
   })
   if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Could not generate refresh token"})
		return
   }

	writeJSONResponse(w, http.StatusOK, user_response{
		Id: new_user.ID.String(),
		Email: new_user.Email,
		Created_at: new_user.CreatedAt.String(),
		Updated_at: new_user.UpdatedAt.String(),
		Token: jwt,
		RefreshToken: refresh_token,
	})
}