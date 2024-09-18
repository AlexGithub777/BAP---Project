package app

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// HandleGetDashboard serves the dashboard page
func (a *App) HandleGetDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println("User Name: ", claims["username"], "User ID: ", claims["user_id"], "User Role: ", claims["role"])

	return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
		"username": claims["username"],
		"role":     claims["role"],
	})
}

// HandleGetAdmin serves the admin page
func (a *App) HandleGetAdmin(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println("User Name: ", claims["username"], "User ID: ", claims["user_id"], "User Role: ", claims["role"])
	return c.Render(http.StatusOK, "admin.html", map[string]interface{}{
		"username": claims["username"],
		"role":     claims["role"],
	})
}

// HandleGetRegister serves the register page
func (a *App) HandleGetRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

// HandleGetForgotPassword serves the forgot password page
func (a *App) HandleGetForgotPassword(c echo.Context) error {
	return c.Render(http.StatusOK, "forgot_password.html", nil)
}
