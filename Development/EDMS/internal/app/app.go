package app

import (
	"log"
	"os"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/config"
	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/database"
	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App holds the application state including database and router
type App struct {
	DB     *database.DB
	Router *echo.Echo
	Logger *log.Logger
}

// handleError is a method of App for handling errors
func (a *App) handleError(c echo.Context, statusCode int, message string, err error) error {
	a.Logger.Printf("\033[31mError: %v\033[0m", err) // Use the logger in the App struct
	return c.JSON(statusCode, map[string]string{"error": message})
}

func (a *App) handleLogger(message string) {
	a.Logger.Printf("\033[34m%s\033[0m", message)
}

// NewApp creates a new instance of App
func NewApp(cfg config.Config) *App {
	// Initialize Echo
	router := echo.New()

	// Set up renderer
	renderer, err := utils.NewTemplateRenderer()
	if err != nil {
		// Handle the error, e.g.:
		panic(err)
	}

	router.Renderer = renderer

	// Serve static files
	router.Static("/static", "static")

	router.Use(middleware.Logger())  // Log requests
	router.Use(middleware.Recover()) // Recover from panics
	router.Use(middleware.CORS())    // Enable CORS

	// Initialize Database
	db, err := database.NewDB(cfg)
	if err != nil {
		panic(err)
	}

	tempFilePath := "internal/seed_complete"
	// Check data import status
	_, err = os.Stat(tempFilePath)
	if os.IsNotExist(err) {
		log.Println("Data not imported")
		database.SeedData(db.DB) // Seed data
	} else {
		log.Println("Data already imported")
	}

	// Initialize Logger
	logger := log.New(os.Stdout, "\033[34mAPP: \033[0m", log.LstdFlags)

	app := &App{
		DB:     db,
		Router: router,
		Logger: logger,
	}

	// Initialize routes
	app.initRoutes()

	return app
}
