package main

import "net/http"

func (cfg *apiConfig) handlerMetricsReset(w http.ResponseWriter, r *http.Request) {
    if cfg.platform != "dev" {
        writeJSONResponse(w, http.StatusForbidden, errResponse{Error: "This can only run on dev"})
        return
    }

    cfg.fileserverHits.Store(0)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Counter reset to 0"))
    err := cfg.dbQueries.DeleteAllUsers(r.Context())   
    if err != nil {
        writeJSONResponse(w, http.StatusBadRequest, errResponse{Error: "Failed to delete users"})
    }
}