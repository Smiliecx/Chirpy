package main

import (
	"encoding/json"
	"net/http"

	"github.com/Smiliecx/Chirpy/internal/auth"
	"github.com/Smiliecx/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChangeData(w http.ResponseWriter, r *http.Request) {
	type  ChangeDataRequest struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type ChangeDataResponse struct {
		Email string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")
   	defer r.Body.Close()

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find bearer token on header"})
		return
	}

	_, err = auth.ValidateJWT(bearer, cfg.authSecret)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Bearer token invalid"})
		return
	}

	var changedData ChangeDataRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&changedData)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Something went wrong"})
		return
   	}

	hashedPassword, err := auth.HashPassword(changedData.Password)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Could not hash password"})
		return
	}

	err = cfg.dbQueries.UpdateUserInfo(r.Context(), database.UpdateUserInfoParams{
		Email: changedData.Email,
		HashedPassword: hashedPassword,
		Token: bearer,
	})
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Had trouble updating the account info"})
		return
	}

	writeJSONResponse(w, http.StatusOK, ChangeDataResponse{
		Email: changedData.Email,
	})
 }
 