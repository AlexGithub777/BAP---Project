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

// GetAllDevices fetches all emergency devices from the database with optional filtering by building code
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
		return c.Redirect(http.StatusSeeOther, "/dashboard?error="+err.Error())
	}

	// Insert new emergency device
	err = a.DB.AddEmergencyDevice(emergencyDevice)
	if err != nil {
		a.handleLogger("Error adding device: " + err.Error())
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Error adding device")
	}

	// Redirect to dashboard with success message
	return c.Redirect(http.StatusFound, "/dashboard?message=Device added successfully")
}

func validateDevice(roomIDStr, emergencyDeviceTypeIDStr, extinguisherTypeIDStr, serialNumber, manufactureDateStr, lastInspectionDateStr, size, description, status string) (*models.EmergencyDevice, error) {
	const (
		ErrDeviceTypeRequired          string = "device Type is required"
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

/*
func (a *App) HandleAddDevice(c echo.Context) error {
    // Parse form data
    roomStr := c.FormValue("room")
    emergencyDeviceTypeIDStr := c.FormValue("emergency_device_type_id")
    serialNumber := c.FormValue("serial_number")
    manufactureDateStr := c.FormValue("manufacture_date")
    lastInspectionDateStr := c.FormValue("last_inspection_date")
    size := c.FormValue("size")
    description := c.FormValue("description")
    status := c.FormValue("status")

    // Validate input
    if roomStr == "" || emergencyDeviceTypeIDStr == "" || serialNumber == "" ||
        manufactureDateStr == "" || size == "" || status == "" {
        return a.handleError(c, http.StatusBadRequest, "All fields are required", nil)
    }

    // Convert room ID and emergency device type ID to integers
    roomID, err := strconv.Atoi(roomStr)
    if err != nil {
        log.Printf("Error converting room to integer: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid room ID", err)
    }

    emergencyDeviceTypeID, err := strconv.Atoi(emergencyDeviceTypeIDStr)
    if err != nil {
        log.Printf("Error converting emergency device type ID to integer: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid emergency device type ID", err)
    }

    // Parse date strings into time.Time format
    manufactureDate, err := time.Parse("2006-01-02", manufactureDateStr)
    if err != nil {
        log.Printf("Error parsing manufacture date: %v", err)
        return a.handleError(c, http.StatusBadRequest, "Invalid manufacture date format", err)
    }

    // Optional: Parse last inspection date if provided
    var lastInspectionDate sql.NullTime
    if lastInspectionDateStr != "" {
        parsedDate, err := time.Parse("2006-01-02", lastInspectionDateStr)
        if err != nil {
            return a.handleError(c, http.StatusBadRequest, "Invalid last inspection date format", err)
        }
        lastInspectionDate = sql.NullTime{Time: parsedDate, Valid: true}
    }

    // Insert new emergency device
    var emergencyDeviceID int
    err = a.DB.QueryRow(`
        INSERT INTO emergency_devices (
            emergency_device_type_id,
            room_id,
            manufacture_date,
            serial_number,
            description,
            size,
            last_inspection_date,
            status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) RETURNING emergency_device_id
    `,
        emergencyDeviceTypeID,
        roomID,
        manufactureDate,
        serialNumber,
        description,
        size,
        lastInspectionDate,
        status).Scan(&emergencyDeviceID)
    if err != nil {
        return a.handleError(c, http.StatusInternalServerError, "Error creating emergency device", err)
    }

    // Create the new EmergencyDevice model
    newDevice := models.EmergencyDevice{
        EmergencyDeviceID:    emergencyDeviceID,
        EmergencyDeviceTypeID: emergencyDeviceTypeID,
        RoomID:               roomID,
        ManufactureDate:      manufactureDate,
        SerialNumber:         serialNumber,
        Description:          description,
        Size:                 size,
        LastInspectionDate:   &lastInspectionDate.Time, // only set if valid
        Status:               status,
    }

    // Build HTML for the new row
    newRowHTML := fmt.Sprintf(`
        <tr>
            <td>%d</td>
            <td>%d</td>
            <td>%d</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
            <td>%s</td>
        </tr>`,
        newDevice.EmergencyDeviceID,
        newDevice.EmergencyDeviceTypeID,
        newDevice.RoomID,
        newDevice.SerialNumber,
        newDevice.ManufactureDate.Format("02-01-2006"),
        newDevice.Size,
        newDevice.Description,
        newDevice.LastInspectionDate.Format("02-01-2006"), // ensure this is set correctly
        newDevice.Status,
    )

    // Return success message and the new row HTML
    return c.JSON(http.StatusOK, map[string]string{
        "message": "Emergency device created successfully.",
        "rowHTML": newRowHTML,
    })
}
*/
