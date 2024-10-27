package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// HandleGetInspectionsByDeviceID fetches all inspections by device ID
func (a *App) HandleGetAllInspectionsByDeviceID(c echo.Context) error {
	deviceID, err := strconv.Atoi(c.QueryParam("device_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	inspections, err := a.DB.GetAllInspectionsByDeviceID(deviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, inspections)
}

// HandleGetInspectionByID fetches an inspection by ID
func (a *App) HandleGetInspectionByID(c echo.Context) error {
	inspectionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	inspection, err := a.DB.GetInspectionByID(inspectionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, inspection)
}
