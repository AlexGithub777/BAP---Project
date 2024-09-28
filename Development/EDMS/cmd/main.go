package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/app"
	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/config"
	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the app
	application := app.NewApp(cfg)

	// Get the local IP that has Internet connectivity
	ip := utils.GetLocalIP().String()

	log.Printf("Starting HTTP service on http://%s:3000", ip)

	// HTTP listener is in a goroutine as its blocking
	go func() {
		if err := application.Router.StartTLS(ip+":3000", "./certs/cert.pem", "./certs/key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Setup a ctrl-c trap to ensure a graceful shutdown
	// This would also allow shutting down other pipes/connections. eg DB
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Log the shutdown process
	log.Println("Shutting HTTP service down")
	if err := application.Router.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}

	log.Println("Closing database connections")
	if err := application.DB.Close(); err != nil {
		log.Fatalf("Database Close Failed: %v", err)
	}

	log.Println("Shutdown complete")
}
