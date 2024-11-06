package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/labstack/echo/v4"
)

func (a *App) HandleGetAllRooms(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	buildingId := c.QueryParam("buildingId")

	rooms, err := a.DB.GetAllRooms(buildingId)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, rooms)
}

func (a *App) HandleGetRoomByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/dashboard?error=Method not allowed",
		})
	}

	roomId := c.Param("id")
	if roomId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Room ID is required",
			"redirectURL": "/dashboard?error=Room ID is required",
		})
	}

	roomIdInt, err := strconv.Atoi(roomId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid room ID",
			"redirectURL": "/dashboard?error=Invalid room ID",
		})
	}

	room, err := a.DB.GetRoomByID(roomIdInt)
	if err != nil {
		a.handleError(c, http.StatusNotFound, "Room not found", err)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Room not found",
			"redirectURL": "/dashboard?error=Room not found",
		})
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, room)
}

func (a *App) HandlePostRoom(c echo.Context) error {
	//Check if request is not a post request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	//Pass room name from form
	roomCode := c.FormValue("addRoomCode")
	a.handleLogger("Room Code: " + roomCode)

	// Parse the building ID
	buildingId := c.FormValue("addRoomBuildingCode")
	if buildingId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building ID is required")
	}

	//Validate room name length
	if len(roomCode) < 1 || len(roomCode) > 100 {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Room Name must be between 1 and 100 characters long")
	}

	// Convert building ID to int
	buildingIdInt, err := strconv.Atoi(buildingId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid building ID")
	}

	// Check if the building exists
	building, err := a.DB.GetBuildingById(buildingIdInt)
	if err != nil || building == nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building does not exist")
	}

	_, err = a.DB.GetRoomByCodeAndSite(roomCode, building.SiteID)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Room already exists at this site")
	}

	// Check if the room already exists at the building
	_, err = a.DB.GetRoomByCodeAndBuilding(roomCode, buildingIdInt)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Room already exists at this building")
	}

	var room = models.Room{
		RoomCode:   roomCode,
		BuildingID: buildingIdInt,
	}

	// Add the room to the database
	err = a.DB.AddRoom(&room)
	if err != nil {
		a.handleLogger("Error adding Room: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error adding room")
	}

	//Redirect to admin settings with sucsess message
	return c.Redirect(http.StatusFound, "/admin?message=Room added sucsessfully")
}

func (a *App) HandlePutRoom(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/admin?error=Method not allowed",
		})
	}

	// Get the room ID from the URL
	roomId := c.Param("id")
	if roomId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Room ID is required",
			"redirectURL": "/admin?error=Room ID is required",
		})
	}

	// Convert the room ID to an int
	roomIdInt, err := strconv.Atoi(roomId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid room ID",
			"redirectURL": "/admin?error=Invalid room ID",
		})
	}

	// Check if the room exists
	_, err = a.DB.GetRoomByID(roomIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Room does not exist",
			"redirectURL": "/admin?error=Room does not exist",
		})
	}

	// Parse the form data from the request body
	var roomDto models.RoomDto
	if err := c.Bind(&roomDto); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid request data",
			"redirectURL": "/admin?error=Invalid request data",
		})
	}

	// Validate the room name
	if len(roomDto.RoomCode) < 1 || len(roomDto.RoomCode) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Room name must be between 1 and 100 characters long",
			"redirectURL": "/admin?error=Room name must be between 1 and 100 characters long",
		})
	}

	// Convert the building ID to an int
	buildingIdInt, err := strconv.Atoi(roomDto.BuildingID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid building ID",
			"redirectURL": "/admin?error=Invalid building ID",
		})
	}

	// Check if the building exists and get the SiteID
	building, err := a.DB.GetBuildingById(buildingIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Building does not exist",
			"redirectURL": "/admin?error=Building does not exist",
		})
	}

	// Check if the room already exists at the site
	existingRoomAtSite, err := a.DB.GetRoomByCodeAndSite(roomDto.RoomCode, building.SiteID)
	if err == nil || existingRoomAtSite != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error":       "Room already exists at this site",
			"redirectURL": "/admin?error=Room already exists at this site",
		})
	}

	// Check if the room already exists at the building
	existingRoom, err := a.DB.GetRoomByCodeAndBuilding(roomDto.RoomCode, buildingIdInt)
	if err == nil || existingRoom != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error":       "Room already exists at this building",
			"redirectURL": "/admin?error=Room already exists at this building",
		})
	}

	// Create a new room object
	room := models.Room{
		RoomID:     roomIdInt,
		RoomCode:   roomDto.RoomCode,
		BuildingID: buildingIdInt,
	}

	// Update the room in the database
	err = a.DB.UpdateRoom(&room)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error updating room",
			"redirectURL": "/admin?error=Error updating room",
		})
	}

	// Redirect to admin settings with success message
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Room updated successfully",
		"redirectURL": "/admin?message=Room updated successfully",
	})
}

func (a *App) HandleDeleteRoom(c echo.Context) error {
	// Check if request is not a DELETE request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/admin?error=Method not allowed",
		})
	}

	// Get the room ID from the URL
	roomId := c.Param("id")
	if roomId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Room ID is required",
			"redirectURL": "/admin?error=Room ID is required",
		})
	}

	// Convert the room ID to an int
	roomIdInt, err := strconv.Atoi(roomId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid room ID",
			"redirectURL": "/admin?error=Invalid room ID",
		})
	}

	// Check if the room exists
	_, err = a.DB.GetRoomByID(roomIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":       "Room does not exist",
			"redirectURL": "/admin?error=Room does not exist",
		})
	}

	// Handle foreign key constraints, check if room has any devices
	devices, err := a.DB.GetDevicesByRoomID(roomIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching devices",
			"redirectURL": "/admin?error=Error fetching devices",
		})
	}

	if len(devices) > 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"error":       "Can't delete room with associated devices",
			"redirectURL": "/admin?error=Can't delete room with associated devices",
		})
	}

	// Delete the room from the database
	err = a.DB.DeleteRoom(roomIdInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting room",
			"redirectURL": "/admin?error=Error deleting room",
		})
	}

	// Redirect to admin settings with success message
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Room deleted successfully",
		"redirectURL": "/admin?message=Room deleted successfully",
	})
}
