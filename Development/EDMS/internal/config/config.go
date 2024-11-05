package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser        string
	DBPassword    string
	DBName        string
	DBHost        string
	DBPort        int
	AdminPassword string
	JWTSecret     string
}

func LoadConfig() Config {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", cwd)

	// Load the .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the DB_PORT environment variable
	dbPortStr := os.Getenv("DB_PORT")

	// Convert the string to an integer and handle any potential errors
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	return Config{
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        dbPort,
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}
}
