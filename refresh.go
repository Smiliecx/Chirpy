package main

import (
	"net/http"
	"time"

	"github.com/Smiliecx/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type refreshResponse struct {
		Token string `json:"token"`
	}

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find bearer token on header"})
		return
	}

	token_row, err := cfg.dbQueries.GetRefreshToken(r.Context(), bearer)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find Refresh token"})
		return
	}

	if token_row.ExpiresAt.Before(time.Now()) {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Token is Expired"})
		return
	}

	if token_row.RevokedAt.Valid {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Token is Revoked"})
		return
	}

	user_id, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token_row.Token)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find user from Token"})
		return
	}

	new_token, _ := auth.MakeJWT(user_id, cfg.authSecret, time.Duration(time.Hour))
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Could not generate jwt"})
		return
   	}

	writeJSONResponse(w, http.StatusOK, refreshResponse{
		Token: new_token,
	})
}