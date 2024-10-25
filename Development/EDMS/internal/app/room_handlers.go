package app

import (
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

	// Check if the building exists
	_, err := a.DB.GetBuildingById(buildingId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building does not exist")
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
