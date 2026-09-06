package main

import (
	"encoding/json"
	"fmt"
	"github.com/haaguileraa/chirpy/internal/database"
	"net/http"
)

const maxBodyLength = 140

func (cfg *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	var chirpReq chirpyChirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode chirp with user id: %v", err))
		return
	}

	validatedChirpReq, err := validateChirp(chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := database.CreateChirpParams {
		Body:	validatedChirpReq.Body, 
		UserID:	validatedChirpReq.UserID,
	}

	chirpDb, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp := chirpDB2JSONChirp(chirpDb)	
	respondWithJSON(w, http.StatusCreated, chirp)
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
		UserID:	chirp.UserID,
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
