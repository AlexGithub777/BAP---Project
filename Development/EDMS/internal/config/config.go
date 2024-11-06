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
	// Try to load .env file, but don't fail if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, will use environment variables")
	}

	// List of required environment variables
	requiredEnvVars := map[string]string{
		"DB_USER":        "",
		"DB_PASSWORD":    "",
		"DB_NAME":        "",
		"DB_HOST":        "",
		"DB_PORT":        "",
		"ADMIN_PASSWORD": "",
		"JWT_SECRET":     "",
	}

	// Check for missing environment variables
	var missingVars []string
	for env := range requiredEnvVars {
		if value := os.Getenv(env); value == "" {
			missingVars = append(missingVars, env)
		}
	}

	// If any required variables are missing, log them and exit
	if len(missingVars) > 0 {
		log.Fatalf("Missing required environment variables: %v", missingVars)
	}

	// Get and validate DB_PORT
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}

	// Create and return the config
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
