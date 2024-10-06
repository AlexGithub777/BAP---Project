package app

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAllDeviceTypes(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	emergencyDeviceTypes, err := a.DB.GetAllDeviceTypes()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, emergencyDeviceTypes)
}

func (a *App) HandlePostDeviceType(c echo.Context) error {
//Check if request is not a post request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	//Pass device type name from form
	deviceTypeName := c.FormValue("device_type_name")
	a.handleLogger("Device Type Name: " + deviceTypeName)

	//Validate device type name
    deviceNameRegex := regexp.MustCompile("^[a-zA-Z0-9_ ]{1,50}$")
    if !deviceNameRegex.MatchString(deviceTypeName) {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error":       "Device Type Name must be between 1 and 50 characters long, and contain only letters, numbers, and underscores.",
            "redirectURL": "/admin?error=Device Type Name must be between 1 and 50 characters long, and contain only letters, numbers, and underscores",
        })
    }
	return nil
}