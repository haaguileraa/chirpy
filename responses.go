package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const maxBodyLength = 140

type chirpyBody struct {
	Body 	string `json:"body"`
}

type chirpyError struct {
	Error 	string `json:"error"`
}

type chirpyCleanedBody struct {
	CleanedBody	string `json:"cleaned_body"`
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	var body chirpyBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("could not decode Chirp: %v", err))
		return
	}
	
	bodyIsInvalid := len(body.Body) > maxBodyLength  

	if bodyIsInvalid {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := getBadWords()
	
	payload := chirpyCleanedBody {
		CleanedBody:	replaceBadWords(body.Body, badWordReplacement, badWords),
	}
	respondWithJSON(w, http.StatusOK, payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("content-type", "application/json")
	chirpyErr := chirpyError {
		Error:	msg,
	}

	data, err := json.Marshal(chirpyErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("could not marshal error data: %v", err)))
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
