package main

import (
	"encoding/json"
	"fmt"
	"github.com/haaguileraa/chirpy/internal/database"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db		*database.Queries
	platform	string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {	
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) handlerNumberRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	metrics := getFormattedMetrics(int(cfg.fileserverHits.Load()))
	w.Write([]byte(metrics))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, http.StatusText(http.StatusOK))
		return
	}
	err := cfg.db.CleanDatabase(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset"))
}

func (cfg *apiConfig) handlerUser(w http.ResponseWriter, r *http.Request) {
	var email chirpyEmail
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode email: %v", err))
		return
	}
	
	userDb, err := cfg.db.CreateUser(r.Context(), email.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	user := User {
		ID: 		userDb.ID,
		CreatedAt:	userDb.CreatedAt,
		UpdatedAt:	userDb.UpdatedAt,
		Email:		userDb.Email,
	}
	respondWithJSON(w, http.StatusCreated, user)
}

