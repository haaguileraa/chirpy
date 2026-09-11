package main

import (
	"encoding/json"
	"fmt"
	"github.com/haaguileraa/chirpy/internal/auth"
	"github.com/haaguileraa/chirpy/internal/database"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db		*database.Queries
	platform	string
	jwtSecret	string
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
	var userReq chirpyUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode user: %v", err))
		return
	}
	
	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	params := database.CreateUserParams {
		Email:	userReq.Email,
		HashedPassword:	hashedPassword,
	}	
	userDb, err := cfg.db.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	user := userDB2JSONUser(userDb) 
	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("error getting token from request"))
		return
	}
	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("error validating jwt: %v", err))
		return
	}
	var userReq chirpyUser
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode new user information: %v", err))
		return
	}
	
	hashedPassword, err := auth.HashPassword(userReq.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	
	editParams := database.UpdateUserPasswordAndEmailParams {
		ID:		userID,
		Email:		userReq.Email,
		HashedPassword:	hashedPassword,
	}		
	
	editedUser, err := cfg.db.UpdateUserPasswordAndEmail(r.Context(), editParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not edit user")
		return
	}
	
	user := userDB2JSONUser(editedUser) 
	respondWithJSON(w, http.StatusOK, user)

}

