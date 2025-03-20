package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Smiliecx/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
	authSecret string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	authSecret := os.Getenv("AUTH_SECRET_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		dbQueries: dbQueries,
		platform: os.Getenv("PLATFORM"), 
		authSecret: authSecret,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirp)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerMetricsReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,    
	}

	log.Fatal(server.ListenAndServe())
}




