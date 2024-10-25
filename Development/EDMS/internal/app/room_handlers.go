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
	roomName := c.FormValue("room_name")
	a.handleLogger("Room Name: " + roomName)

	// Parse the building ID
	buildingId := c.FormValue("building_id")
	if buildingId == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building ID is required")
	}

	// Check if the building exists
	_, err := a.DB.GetBuildingById(buildingId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Building does not exist")
	}

	//Validate room name length
	if len(roomName) < 1 || len(roomName) > 100 {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Room Name must be between 1 and 100 characters long")
	}

	// Convert building ID to int
	buildingIdInt, err := strconv.Atoi(buildingId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid building ID")
	}

	var room = models.Room{
		RoomCode:   roomName,
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
