package app

import (
	"log"
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

func (a *App) HandleGetBuildingByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	buildingId := c.Param("id")
	building, err := a.DB.GetBuildingById(buildingId)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, building)
}

func (a *App) HandlePostBuilding(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form data
	siteId := c.FormValue("addBuildingSite")
	buildingCode := strings.TrimSpace(c.FormValue("addBuildingCode")) // Trim whitespace

	// Validate input
	if siteId == "" || buildingCode == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=All fields are required")
	}

	// Validate site ID
	_, err := a.DB.GetSiteByID(siteId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid site ID")
	}

	// Validate building code length doesnt exceed 100 characters
	if len(buildingCode) > 100 {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building code too long, must be less than 100 characters")
	}

	// Save site information and file path in the database
	siteIdNum, err := strconv.Atoi(siteId)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error converting site ID", err)
	}

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
	// Check if request is a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the id from the URL parameter
	buildingID := c.Param("id")

	// Parse the form data
	var building models.BuildingDto
	if err := c.Bind(&building); err != nil {
		a.handleLogger("Error parsing form data")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Error parsing form data",
			"redirectURL": "/admin?error=Error parsing form data",
		})
	}

	building.BuildingID = buildingID

	log.Printf("Building ID: %s, Building Code: %s, Site ID: %s", buildingID, building.BuildingCode, building.SiteID)

	// Validate input
	if building.SiteID == "" || building.BuildingCode == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "All fields are required",
			"redirectURL": "/admin?error=All fields are required",
		})
	}

	// Validate site ID
	_, err := a.DB.GetSiteByID(building.SiteID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid site ID",
			"redirectURL": "/admin?error=Invalid site ID",
		})
	}

	// Validate building code length doesnt exceed 100 characters
	if len(building.BuildingCode) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Building code too long, must be less than 100 characters",
			"redirectURL": "/admin?error=Building code too long, must be less than 100 characters",
		})
	}

	// Save site information and file path in the database
	siteIdNum, err := strconv.Atoi(building.SiteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error converting site ID",
			"redirectURL": "/admin?error=Error converting site ID",
		})
	}

	buildingIDNum, err := strconv.Atoi(buildingID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error converting building ID",
			"redirectURL": "/admin?error=Error converting building ID",
		})
	}

	buildingModel := &models.Building{
		BuildingID:   buildingIDNum,
		SiteID:       siteIdNum,
		BuildingCode: building.BuildingCode,
	}

	err = a.DB.UpdateBuilding(buildingModel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error updating building",
			"redirectURL": "/admin?error=Error updating building",
		})
	}

	// Respond to the client
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Building updated successfully",
		"redirectURL": "/admin?message=Building updated successfully",
	})
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
