package main

import (
	_ "github.com/lib/pq"

	"database/sql"
	"github.com/haaguileraa/chirpy/internal/database"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	cfg := apiConfig {
		db: dbQueries,
		platform: platform,
	}
	serve(&cfg)
}
