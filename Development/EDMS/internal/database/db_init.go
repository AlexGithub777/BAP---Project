package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AlexGithub777/safety-device-app/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
	*sql.DB
}

func NewDB(cfg config.Config) (*DB, error) {
	log.Println("Connecting to database...")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, cfg.DBPort)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connection is established
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	return &DB{db}, nil
}
