package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/haaguileraa/chirpy/internal/database"
	"github.com/haaguileraa/chirpy/internal/auth"
	"net/http"
)

const maxBodyLength = 140

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error getting token from request"))
		return
	}
	var chirpReq chirpyChirp
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode chirp with user id: %v", err))
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Error validating jwt: %s", err.Error()))
		return
	}

	validatedChirpReq, err := validateChirp(chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := database.CreateChirpParams {
		Body:	validatedChirpReq.Body, 
		UserID:	userID,
	}

	chirpDb, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	chirp := chirpDB2JSONChirp(chirpDb)	
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirpsDb, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirps := make([]Chirp, len(chirpsDb))

	for i, chirpDb := range chirpsDb {
		chirps[i] = chirpDB2JSONChirp(chirpDb)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirpDb, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
	}
	
	chirp := chirpDB2JSONChirp(chirpDb)
	respondWithJSON(w, http.StatusOK, chirp)
}


func validateChirp(chirp chirpyChirp) (chirpyChirp, error) {
	bodyIsInvalid := len(chirp.Body) > maxBodyLength  

	if bodyIsInvalid {
		return chirpyChirp{}, fmt.Errorf("Chirp is too long")
	}
	badWords := getBadWords()

	cleanedBody := replaceBadWords(chirp.Body, badWordReplacement, badWords)
	return chirpyChirp {
		Body:	cleanedBody,
	}, nil
}

func chirpDB2JSONChirp(chirpDb database.Chirp) Chirp {
	return Chirp {
		ID:		chirpDb.ID,
		CreatedAt:	chirpDb.CreatedAt,
		UpdatedAt:	chirpDb.UpdatedAt,
		Body:		chirpDb.Body,
		UserID:		chirpDb.UserID,
	}
}
