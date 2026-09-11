package main

import (
	"fmt"
	"github.com/haaguileraa/chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request, refreshToken string) {	
	userID, err := cfg.db.GetUserIDFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("could not get user from refresh token: %v", err))
		return
	}
	jwtAccessExpiresIn := defaultExpirationAccessSeconds * time.Second
	token, err := auth.MakeJWT(userID, cfg.jwtSecret, jwtAccessExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create")
		return
	}
	
	type tokenJson struct {
		Token	string	`json:"token"`
	}

	response := tokenJson {
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func getBearerTokenMiddleware(next func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("error getting token from request"))
			return
		}
		next(w, r, token)
	}
}

