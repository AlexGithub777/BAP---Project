package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func (a *App) initRoutes() {
	// Public routes
	a.Router.GET("/", a.HandleGetLogin)
	a.Router.GET("/register", a.HandleGetRegister)
	a.Router.GET("/forgot-password", a.HandleGetForgotPassword)
	a.Router.POST("/forgot-password", a.HandlePostForgotPassword)
	a.Router.POST("/register", a.HandlePostRegister)
	a.Router.POST("/login", a.HandlePostLogin)
	a.Router.GET("/logout", a.HandleGetLogout)

	// Protected routes
	protected := a.Router.Group("")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("JWT_SECRET")),
		TokenLookup: "cookie:token",
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			fmt.Println("User Name: ", claims["username"], "User ID: ", claims["user_id"], "User Role: ", claims["role"])
		},
		ErrorHandler: func(c echo.Context, err error) error {
			return c.Redirect(http.StatusSeeOther, "/")
		},
	}))

	protected.GET("/dashboard", a.HandleGetDashboard)
	protected.GET("/admin", a.HandleGetAdmin)
	// Add the rest of your protected routes here

	// API routes
	api := protected.Group("/api")
	api.GET("/user", a.HandleGetAllUsers)
	api.GET("/emergency-device", a.HandleGetAllDevices)
	api.GET("/emergency-device-type", a.HandleGetAllDeviceTypes)
	api.GET("/extinguisher-type", a.HandleGetAllExtinguisherTypes)
	api.GET("/room", a.HandleGetAllRooms)
	api.GET("/building", a.HandleGetAllBuildings)
	api.GET("/site", a.HandleGetAllSites)
	//a.Router.POST("/api/emergency-device", a.HandleAddDevice)
	// Add the rest of your API routes here
}
