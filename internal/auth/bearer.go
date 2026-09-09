package auth

import (
	"fmt"
	"net/http"
	"strings"
)

const malformattedHeader = "Malformatted Authorization header"

func GetBearerToken(headers http.Header) (string, error) {
	bearer, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("could not find Authorization header")
	}

	prefix := "bearer "
	// bearer is case insensitive -> ToLower
	token, found := strings.CutPrefix(strings.TrimSpace(strings.ToLower(bearer[0])), prefix)
	if !found {
		return "", fmt.Errorf("%s. Expecting 'Bearer <TOKEN>' got '%s'", malformattedHeader, bearer)
	}

	trimmedToken := strings.TrimSpace(token)
	if trimmedToken == "" {
		return "", fmt.Errorf("%s. No token was found.", malformattedHeader)
	}
	return trimmedToken, nil
}
