package main 

import (
	"encoding/json"
	"fmt"
	"github.com/haaguileraa/chirpy/internal/auth"
	"github.com/haaguileraa/chirpy/internal/database"
	"net/http"
	"time"
)

const (
	defaultExpirationAccessSeconds = 3600
	defaultExpirationRefreshHours = 60 * 24
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var userReq chirpyUser
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode user: %v", err))
		return
	}

	userDb, err := cfg.db.GetUserByEmail(r.Context(), userReq.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not obtain user: %v", err))
		return
	}
	
	match, err := auth.CheckPasswordHash(userReq.Password, userDb.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not check password")
		return
	}

	if !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	
	jwtAccessExpiresIn := defaultExpirationAccessSeconds * time.Second
	token, err := auth.MakeJWT(userDb.ID, cfg.jwtSecret, jwtAccessExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not log in")
		return
	}
	
	paramsRefreshToken := database.CreateRefreshTokenParams {
		Token: 		auth.MakeRefreshToken(),
		UserID: 	userDb.ID,
		ExpiresAt:	time.Now().UTC().Add(defaultExpirationRefreshHours * time.Hour),
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), paramsRefreshToken) 
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token")
		return
	}

	user := userDB2JSONUser(userDb)
	userWithToken := UserWithToken {
		User:		user,
		Token:		token,
		RefreshToken:	refreshToken.Token,
	}
	respondWithJSON(w, http.StatusOK, userWithToken)
}



func userDB2JSONUser(userDb database.User) User {
	return User {
		ID: 		userDb.ID,
		CreatedAt:	userDb.CreatedAt,
		UpdatedAt:	userDb.UpdatedAt,
		Email:		userDb.Email,
	}
}
