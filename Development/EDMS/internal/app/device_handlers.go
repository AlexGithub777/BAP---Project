package app

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/labstack/echo/v4"
)

// HandleGetAllDevices fetches all emergency devices from the database with optional filtering by building code
// and returns the results as JSON
func (a *App) HandleGetAllDevices(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}
	siteId := c.QueryParam("site_id")
	buildingCode := c.QueryParam("building_code")

	emergencyDevices, err := a.DB.GetAllDevices(siteId, buildingCode)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, emergencyDevices)
}

// HandleGetDeviceByID fetches a single emergency device by ID from the database and returns the result as JSON
func (a *App) HandleGetDeviceByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	// Get the device ID from the URL
	deviceIDStr := c.Param("id")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		return a.handleError(c, http.StatusBadRequest, "Invalid device ID", err)
	}

	// Fetch the device from the database
	device, err := a.DB.GetDeviceByID(deviceID)

	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the result as JSON
	return c.JSON(http.StatusOK, device)
}

func (a *App) HandlePostDevice(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}
	// Parse form data
	roomIDStr := c.FormValue("room_id")
	emergencyDeviceTypeIDStr := c.FormValue("emergency_device_type")
	extinguisherTypeIDStr := c.FormValue("extinguisher_type")
	serialNumber := c.FormValue("serial_number")
	manufactureDateStr := c.FormValue("manufacture_date")
	lastInspectionDateStr := c.FormValue("last_inspection_date")
	size := c.FormValue("size")
	description := c.FormValue("description")
	status := c.FormValue("status")
	a.handleLogger("Room: " + roomIDStr)
	a.handleLogger("Emergency Device Type: " + emergencyDeviceTypeIDStr)
	a.handleLogger("extinguisher_type_id: " + extinguisherTypeIDStr)
	a.handleLogger("serial_number: " + serialNumber)
	a.handleLogger("manufacture_date: " + manufactureDateStr)
	a.handleLogger("last_inspection_date: " + lastInspectionDateStr)
	a.handleLogger("size: " + size)
	a.handleLogger("description: " + description)
	a.handleLogger("status: " + status)

	// Validate input
	emergencyDevice, err := validateDevice(roomIDStr, emergencyDeviceTypeIDStr, extinguisherTypeIDStr, serialNumber, manufactureDateStr, lastInspectionDateStr, size, description, status)
	if err != nil {
		a.handleLogger("Error validating device: " + err.Error())
		// Redirect to dashboard with error message
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+"Error validating device: "+err.Error())
	}

	// Insert new emergency device
	err = a.DB.AddEmergencyDevice(emergencyDevice)
	if err != nil {
		a.handleLogger("Error adding device: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+err.Error())
	}

	// Redirect to dashboard with success message
	return c.Redirect(http.StatusFound, "/dashboard?message=Device added successfully")
}

func (a *App) HandlePutDevice(c echo.Context) error {
	// Check if request is not a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	// Parse the device ID from the URL parameter
	deviceIDStr := c.Param("id")

	// Parse form data from the request body
	var device models.EmergencyDeviceDto
	if err := c.Bind(&device); err != nil {
		a.handleLogger("Error binding request body: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid request body",
			"redirectURL": "/dashboard?error=Invalid request body",
		})
	}

	// Convert the device ID to an integer
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		a.handleLogger("Error converting device ID to integer: " + err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid device ID",
			"redirectURL": "/dashboard?error=Invalid device ID"})
	}

	// Log the incoming data
	a.handleLogger("Device ID: " + deviceIDStr)
	a.handleLogger("Room: " + device.RoomID)
	a.handleLogger("Emergency Device Type: " + device.EmergencyDeviceTypeID)
	a.handleLogger("Extinguisher Type: " + device.ExtinguisherTypeID)
	a.handleLogger("Serial Number: " + device.SerialNumber)
	a.handleLogger("Manufacture Date: " + device.ManufactureDate)
	a.handleLogger("Last Inspection Date: " + device.LastInspectionDate)
	a.handleLogger("Size: " + device.Size)
	a.handleLogger("Description: " + device.Description)
	a.handleLogger("Status: " + device.Status)

	// Validate input
	emergencyDevice, err := validateDevice(device.RoomID, device.EmergencyDeviceTypeID, device.ExtinguisherTypeID, device.SerialNumber, device.ManufactureDate, device.LastInspectionDate, device.Size, device.Description, device.Status)
	if err != nil {
		a.handleLogger("Error validating device: " + err.Error())
		// Redirect to dashboard with error message
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error validating device: " + err.Error(),
			"redirectURL": "/dashboard?error=" + err.Error()})
	}

	// Add the device ID to the emergency device model
	emergencyDevice.EmergencyDeviceID = deviceID

	// Update the device in the database
	err = a.DB.UpdateEmergencyDevice(emergencyDevice)
	if err != nil {
		a.handleLogger("Error updating device: " + err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error updating device: " + err.Error(),
			"redirectURL": "/dashboard?error=" + err.Error()})
	}

	// Redirect to dashboard with success message
	return c.JSON(http.StatusOK, map[string]string{"message": "Device updated successfully", "redirectURL": "/dashboard?message=Device updated successfully"})
}

func validateDevice(roomIDStr, emergencyDeviceTypeIDStr, extinguisherTypeIDStr, serialNumber, manufactureDateStr, lastInspectionDateStr, size, description, status string) (*models.EmergencyDevice, error) {
	const (
		ErrDeviceTypeRequired          string = "device type is required"
		ErrRoomRequired                string = "room is required"
		ErrInvalidRoomID               string = "invalid room ID"
		ErrInvalidEmergencyDeviceID    string = "invalid emergency device type ID"
		ErrInvalidExtinguisherTypeID   string = "invalid extinguisher type ID"
		ErrRoomDoesNotExist            string = "room does not exist"
		ErrDeviceTypeDoesNotExist      string = "emergency Device Type does not exist"
		ErrExtinguisherTypeNotExist    string = "extinguisher Type does not exist"
		ErrInvalidManufactureDate      string = "invalid manufacture date format"
		ErrManufactureDateInFuture     string = "manufacture date cannot be in the future"
		ErrInvalidInspectionDate       string = "invalid last inspection date format"
		ErrInspectionDateInFuture      string = "last inspection date cannot be in the future"
		ErrInspectionBeforeManufacture string = "last inspection date cannot be before manufacture date"
		ErrSerialNumberTooLong         string = "serial number is too long, maximum 50 characters"
		ErrDescriptionTooLong          string = "description is too long, maximum 255 characters"
		ErrSizeTooLong                 string = "size is too long, maximum 50 characters"
		ErrStatusTooLong               string = "status is too long, maximum 50 characters"
	)

	var device models.EmergencyDevice

	if roomIDStr == "" {
		return &device, errors.New(ErrRoomRequired)
	}

	if emergencyDeviceTypeIDStr == "" {
		return &device, errors.New(ErrDeviceTypeRequired)
	}

	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return &device, errors.New(ErrInvalidRoomID)
	}

	emergencyDeviceTypeID, err := strconv.Atoi(emergencyDeviceTypeIDStr)
	if err != nil {
		return &device, errors.New(ErrInvalidEmergencyDeviceID)
	}

	var extinguisherTypeID sql.NullInt64
	if extinguisherTypeIDStr != "" {
		extinguisherTypeID.Int64, err = strconv.ParseInt(extinguisherTypeIDStr, 10, 32)
		if err != nil {
			return &device, errors.New(ErrInvalidExtinguisherTypeID)
		}
		extinguisherTypeID.Valid = true
	} else {
		extinguisherTypeID.Valid = false
	}

	manufactureDate, err := parseDate(manufactureDateStr)
	if err != nil {
		return &device, errors.New(ErrInvalidManufactureDate)
	}

	if manufactureDate.Valid && manufactureDate.Time.After(time.Now()) {
		return &device, errors.New(ErrManufactureDateInFuture)
	}

	lastInspectionDate, err := parseDate(lastInspectionDateStr)
	if err != nil {
		return &device, errors.New(ErrInvalidInspectionDate)
	}

	if lastInspectionDate.Valid && lastInspectionDate.Time.After(time.Now()) {
		return &device, errors.New(ErrInspectionDateInFuture)
	}

	if lastInspectionDate.Valid && manufactureDate.Valid && lastInspectionDate.Time.Before(manufactureDate.Time) {
		return &device, errors.New(ErrInspectionBeforeManufacture)
	}

	if len(serialNumber) > 50 {
		return &device, errors.New(ErrSerialNumberTooLong)
	}

	if len(description) > 255 {
		return &device, errors.New(ErrDescriptionTooLong)
	}

	if len(size) > 50 {
		return &device, errors.New(ErrSizeTooLong)
	}

	if len(status) > 50 {
		return &device, errors.New(ErrStatusTooLong)
	}

	// Set the values of the device model
	// Initialize sql.NullString for optional fields
	device.SerialNumber = sql.NullString{String: serialNumber, Valid: serialNumber != ""}
	device.Size = sql.NullString{String: size, Valid: size != ""}
	device.Description = sql.NullString{String: description, Valid: description != ""}
	device.Status = sql.NullString{String: status, Valid: status != ""}
	device.RoomID = roomID
	device.EmergencyDeviceTypeID = emergencyDeviceTypeID
	device.ExtinguisherTypeID = extinguisherTypeID
	device.ManufactureDate = manufactureDate
	device.LastInspectionDate = lastInspectionDate

	return &device, nil
}

func parseDate(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{}, nil
	}

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{}, err
	}

	return sql.NullTime{Time: parsedDate, Valid: true}, nil
}

func (a *App) HandleDeleteDevice(c echo.Context) error {
	// Check if request is not a DELETE request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}

	// Parse the device ID from the URL parameter
	deviceIDStr := c.Param("id")

	// Convert the device ID to an integer
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid device ID",
			"redirectURL": "/dashboard?error=Invalid device ID",
		})
	}

	// Delete the device from the database
	err = a.DB.DeleteEmergencyDevice(deviceID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting device",
			"redirectURL": "/dashboard?error=Error deleting device: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Device deleted successfully",
		"redirectURL": "/dashboard?message=Device deleted successfully",
	})
}
