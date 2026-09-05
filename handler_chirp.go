package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/haaguileraa/chirpy/internal/database"
	"net/http"
)

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	
	type chirpyChirp struct {
		Body	string    `json:"body"`
		UserID	uuid.UUID `json:"user_id"`
	}

	var chirpReq chirpyChirp
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode chirp with user id: %v", err))
		return
	}
	params := database.CreateChirpParams {
		Body:	chirpReq.Body, 
		UserID:	chirpReq.UserID,
	}

	chirpDb, err := cfg.db.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirp := Chirp {
		ID:		chirpDb.ID,
		CreatedAt:	chirpDb.CreatedAt,
		UpdatedAt:	chirpDb.UpdatedAt,
		Body:		chirpDb.Body,
		UserID:		chirpDb.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}
