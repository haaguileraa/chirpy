package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestJWTPass(t *testing.T) {
	
	userID := uuid.New()
	tokenSecret := "pablitoclavounclavito"
	expiresIn := 100 * time.Second

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("could not make token string, %T", err)
		return
	}
	
	outUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Errorf("wasn't expecting error: %T", err)
		return
	}
	
	if userID != outUserID {
		t.Errorf("wrong user ID")
	}
}

func TestJWTExpired(t *testing.T) {
	
	userID := uuid.New()
	tokenSecret := "pablitoclavounclavito"
	expiresIn := 0 * time.Millisecond

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("could not make token string, %T", err)
		return
	}
	
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("was expecting error")
	}
}

func TestJWTWrongTokenSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "pablitoclavounclavito"
	expiresIn := 100 * time.Second

	tokenString, err := MakeJWT(userID, "enlacabezadeuncalvito", expiresIn)
	if err != nil {
		t.Errorf("could not make token string, %T", err)
		return
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("was expecting error")
	}
}
