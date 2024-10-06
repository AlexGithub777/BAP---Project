package app

import (
	"net/http"
	"regexp"
	"strconv"

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
        return c.Redirect(http.StatusSeeOther, "/admin?error=Device Type Name must be between 1 and 50 characters long, and contain only letters, numbers, and underscores")
    }

	//Check device type name is unique
    if _, err := a.DB.GetDeviceTypeByName(deviceTypeName); err == nil {
        return c.Redirect(http.StatusSeeOther, "/admin?error=Device Type Name already exists")
    }

	err := a.DB.AddEmergencyDeviceType(deviceTypeName)
	if err != nil {
		a.handleLogger("Error adding Device Type: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error adding device type")
	}

	//Redirect to admin settings with sucsess message
	return c.Redirect(http.StatusFound, "/admin?message=Device Type added sucsessfully")
}

func (a *App) HandleGetAllDeviceTypeByID(c echo.Context) error {
	//Check if request is not a get request
	if c.Request().Method != http.MethodGet {
        return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
    }
	deviceTypeIDStr := c.Param("id")
	deviceTypeID, err := strconv.Atoi(deviceTypeIDStr)
	if err != nil {
		a.handleLogger("Invalid Device Type ID")
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Invalid Device Type ID"})
	}

	//Fetch the device type from the database
	deviceType, err := a.DB.GetEmergencyDeviceTypeByID(deviceTypeID)
	if err != nil {
		a.handleLogger("Error fetching Device Type")
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Error fetching Device Type"})
	}

	//return the result as JSON
	return c.JSON(http.StatusOK, deviceType)
}

func (a *App) HandleDeleteDeviceType(c echo.Context) error {
	//Check if request is not a delete request
    if c.Request().Method != http.MethodDelete {
        return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
    }

    // Parse the device type ID from the URL parameter
    emergencyDeviceTypeIDStr := c.Param("id")

    // Convert the device type ID to an integer
    emergencyDeviceTypeID, err := strconv.Atoi(emergencyDeviceTypeIDStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error":       "Invalid Emergency Device Type ID",
            "redirectURL": "/admin?error=Invalid device type ID",
        })
    }

	//get devices by device type
	emergencyDevices, err := a.DB.GetDevicesByTypeID(emergencyDeviceTypeID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":		"Error fetching emergency devices",
			"redirectURL": "/admin?error=Error fetching emergency devices",
		})
	}

	//Check if any devices associated with type to be deleted - display error if so
	if len(emergencyDevices) > 0 {
		return c.JSON(http.StatusSeeOther, map[string]string{
			"error":		"Cannot delete emergency device type with associated devices",
			"redirectURL": "/admin?error=Cannot delete emergency device type with associated devices",
		})
	}

    // Delete the device type from the database
    err = a.DB.DeleteEmergencyDeviceType(emergencyDeviceTypeID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error":       "Error deleting device type",
            "redirectURL": "/admin?error=Error deleting device type: " + err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]string{
        "message":     "Device type deleted successfully",
        "redirectURL": "/admin?message=Device type deleted successfully",
    })
}
