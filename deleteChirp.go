package main

import (
	"net/http"
	"strings"

	"github.com/Smiliecx/Chirpy/internal/auth"
	"github.com/Smiliecx/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	idString := strings.TrimPrefix(r.URL.Path, "/api/chirps/")
	uuid, err := uuid.Parse(idString)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Failed to convert UUID string too UUID"})
		return
	}

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find bearer token on header"})
		return
	}

	user_id, err := auth.ValidateJWT(bearer, cfg.authSecret)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Trouble validating bearer token"})
		return
	}

	wasFound, err := cfg.dbQueries.UserOwnsChirp(r.Context(), database.UserOwnsChirpParams{
		UserID: user_id,
		ID: uuid,
	} )
	if err != nil || !wasFound {
		writeJSONResponse(w, http.StatusForbidden, errResponse{Error: "Not the correct user"})
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), uuid)
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, errResponse{Error: "Couldn't delete chirp"})
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}