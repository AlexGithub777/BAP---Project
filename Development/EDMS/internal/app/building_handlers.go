package app

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAllBuildings(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	siteId := c.QueryParam("siteId")
	buildings, err := a.DB.GetAllBuildings(siteId)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, buildings)
}

func (a *App) HandlePostBuilding(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form, limiting upload size to 10MB
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error parsing form")
	}

	// Get the form values
	siteId := strings.TrimSpace(c.FormValue("addSiteId"))             // Trim whitespace 	// Trim whitespace
	buildingCode := strings.TrimSpace(c.FormValue("addBuildingCode")) // Trim whitespace

	// Validate input
	if siteId == "" || buildingCode == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=All fields are required")
	}

	// Save site information and file path in the database
	siteIdNum, err := strconv.Atoi(siteId)
	building := &models.Building{
		SiteID:       siteIdNum,
		BuildingCode: buildingCode,
	}

	err = a.DB.AddBuilding(building)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error saving building", err)
	}

	// Redirect to the admin page with a success message
	return c.Redirect(http.StatusFound, "/admin?message=Building added successfully")
}

func (a *App) HandleEditBuilding(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form, limiting upload size to 10MB
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error parsing form")
	}

	// Get the form values
	//buildingID := c.FormValue("editBuildingID")
	buildingID := c.Param("id")
	siteID := c.FormValue("editSiteID")
	buildingCode := strings.TrimSpace(c.FormValue("editBuildingCode")) // Trim whitespace

	buildingIdNum, err := strconv.Atoi(buildingID)
	siteIdNum, err := strconv.Atoi(siteID)
	building := &models.Building{
		BuildingID:   buildingIdNum,
		SiteID:       siteIdNum,
		BuildingCode: buildingCode,
	}

	err = a.DB.UpdateBuilding(building)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error saving building", err)
	}

	// Respond to the client
	return c.Redirect(http.StatusFound, "/admin?message=Building updated successfully")
}

func (a *App) HandleDeleteBuilding(c echo.Context) error {
	// Check if request is not a delete request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/admin?error=Method not allowed",
		})
	}

	// Get the site ID from the URL
	buildingID := c.Param("id")

	// Get the building by ID
	building, err := a.DB.GetBuildingById(buildingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching building",
			"redirectURL": "/admin?error=Error fetching building",
		})
	}

	// Check if the site has any emergency devices
	emergencyDevices, err := a.DB.GetAllDevices(strconv.Itoa(building.SiteID), building.BuildingCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching emergency devices",
			"redirectURL": "/admin?error=Error fetching emergency devices",
		})
	}

	// Check if the site has any emergency devices
	if len(emergencyDevices) > 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot delete building with associated emergency devices",
			"redirectURL": "/admin?error=Cannot delete building with associated emergency devices",
		})
	}

	// Check if the building has any rooms
	rooms, err := a.DB.GetRoomsByBuildingID(buildingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching rooms",
			"redirectURL": "/admin?error=Error fetching rooms",
		})
	}

	// Check if the site has any rooms
	if len(rooms) > 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot building site with associated rooms",
			"redirectURL": "/admin?error=Cannot delete building with associated rooms",
		})
	}

	// Delete the building from the database
	err = a.DB.DeleteBuilding(buildingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting site",
			"redirectURL": "/admin?error=Error deleting site",
		})
	}

	// Respond to the client
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Building deleted successfully",
		"redirectURL": "/admin?message=Building deleted successfully",
	})
}
