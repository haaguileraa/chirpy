package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request, refreshToken string) {
	err := cfg.db.SetRevokeAtNowByToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("could not find token: %v", err))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

