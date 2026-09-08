package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const issuer = "chirpy-access"

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims {
		Issuer:		issuer,
		Subject:	userID.String(),
		IssuedAt:	&jwt.NumericDate{time.Now().UTC()},
		ExpiresAt:	&jwt.NumericDate{time.Now().Add(expiresIn).UTC()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	
	claims, ok := token.Claims.(*jwt.RegisteredClaims)

	if !ok {
		return uuid.Nil, fmt.Errorf("could not obtain claims")
	}

	if claims.Issuer != issuer {
		return uuid.Nil, fmt.Errorf("issuer does not match")
	}

	userID, err  := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("could not parse user id")
	}
	
	return userID, nil
}
