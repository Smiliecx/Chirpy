package main

import (
	"net/http"

	"github.com/Smiliecx/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Could not find bearer token on header"})
		return
	}

	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), bearer)
	if err != nil {
		writeJSONResponse(w, http.StatusUnauthorized, errResponse{Error: "Error while attempting to revoke token"})
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}