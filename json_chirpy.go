package main

import (
	"github.com/google/uuid"
	"time"
)

type chirpyChirp struct {
	Body	string    `json:"body"`
}

type chirpyError struct {
	Error 	string `json:"error"`
}

type chirpyUser struct {
	Email			string 	`json:"email"`
	Password		string 	`json:"password"`
	ExpiresInSeconds	int	`json:"expires_in_seconds"`
}

type User struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string	  `json:"email"`
}

type UserWithToken struct {
	User
	Token	string `json:"token"`
}

type Chirp struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Body		string    `json:"body"`
	UserID		uuid.UUID `json:"user_id"`
}
