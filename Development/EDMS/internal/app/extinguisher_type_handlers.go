package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAllExtinguisherTypes(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	extinguisherTypes, err := a.DB.GetAllExtinguisherTypes()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, extinguisherTypes)
}
