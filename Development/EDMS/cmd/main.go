package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/app"
)

func main() {

	// Initialize the app
	application := app.NewApp()

	// Get PORT from environment, default to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting HTTP service on port %s", port)

	// HTTP listener is in a goroutine as it's blocking
	go func() {
		if err := application.Router.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Setup a ctrl-c trap to ensure a graceful shutdown
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

	log.Println("Shutdown complete")
}
