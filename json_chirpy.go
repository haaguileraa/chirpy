package main

import (
	"github.com/google/uuid"
	"time"
)

type chirpyChirp struct {
	Body	string    `json:"body"`
	UserID	uuid.UUID `json:"user_id"`
}

type chirpyError struct {
	Error 	string `json:"error"`
}

type chirpyEmail struct {
	Email	string `json:"email"`
}

type User struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Email		string    `json:"email"`
}

type Chirp struct {
	ID		uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
	Body		string    `json:"body"`
	UserID		uuid.UUID `json:"user_id"`
}
