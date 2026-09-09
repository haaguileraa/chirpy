package auth

import (
	"fmt"
	"net/http"
	"strings"
)

const malformattedHeader = "Malformatted Authorization header"

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("could not find Authorization header")
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", fmt.Errorf("%s. Expecting 'Bearer <TOKEN>' got '%s'", malformattedHeader, authHeader)
	}

	trimmedToken := strings.TrimSpace(splitAuth[1])
	if trimmedToken == "" {
		return "", fmt.Errorf("%s. No token was found.", malformattedHeader)
	}
	return trimmedToken, nil
}
