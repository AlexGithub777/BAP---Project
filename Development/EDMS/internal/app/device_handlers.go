package app

import (
	"database/sql"
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
	//Validating device details
	if roomIDStr == "" || emergencyDeviceTypeIDStr == "" {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Room and Device Type are required")
	}

	// Convert room ID an to integer
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid room ID")
	}

	// Convert emergency device type ID to integer
	emergencyDeviceTypeID, err := strconv.Atoi(emergencyDeviceTypeIDStr)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid emergency device type ID")
	}

	// Convert extinguisher type ID to integer if provided
	var extinguisherTypeID sql.NullInt64
	if extinguisherTypeIDStr != "" {
		var err error
		extinguisherTypeID.Int64, err = strconv.ParseInt(extinguisherTypeIDStr, 10, 32)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid extinguisher type ID")
		}
		extinguisherTypeID.Valid = true
	} else {
		extinguisherTypeID.Valid = false
	}

	//validate room_id exists
	room, err := a.DB.GetRoomByID(roomID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Error fetching room")
	}

	if room.RoomID != roomID {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Room does not exist")
	}

	//validate emergency_device_type_id exists
	emergencyDeviceType, err := a.DB.GetEmergencyDeviceTypeByID(emergencyDeviceTypeID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Emergency Device Type does not exist")
	}

	if emergencyDeviceType.EmergencyDeviceTypeID != emergencyDeviceTypeID {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Emergency Device Type does not exist")
	}

	// Validate extinguisher type ID exists if provided
	if extinguisherTypeID.Valid {
		// convert extinguisher type ID to int
		extinguisherTypeIDInt := int(extinguisherTypeID.Int64)
		extinguisherType, err := a.DB.GetExtinguisherTypeByID(extinguisherTypeIDInt)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/dashboard?error=Extinguisher Type does not exist")
		}

		if extinguisherType.ExtinguisherTypeID != extinguisherTypeIDInt {
			return c.Redirect(http.StatusSeeOther, "/dashboard?error=Extinguisher Type does not exist")
		}
	}

	// Parse date strings into sql.NullTime format
	manufactureDate := sql.NullTime{}

	if manufactureDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", manufactureDateStr)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid manufacture date format")
		}
		manufactureDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	// Check if manufacture date is in the future
	if manufactureDate.Valid && manufactureDate.Time.After(time.Now()) {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Manufacture date cannot be in the future")
	}

	// Parse last inspection date if provided
	lastInspectionDate := sql.NullTime{}

	if lastInspectionDateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", lastInspectionDateStr)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid last inspection date format")
		}
		lastInspectionDate = sql.NullTime{Time: parsedDate, Valid: true}
	}

	// Check if last inspection date is in the future
	if lastInspectionDate.Valid && lastInspectionDate.Time.After(time.Now()) {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Last inspection date cannot be in the future")
	}

	// check serial number length is not more than 50
	if len(serialNumber) > 50 {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Serial number is too long, maximum 50 characters")
	}

	// Check description length is not more than 255 characters
	if len(description) > 255 {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Description is too long, maximum 255 characters")
	}

	// Check size length is not more than 50 characters
	if len(size) > 50 {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Size is too long, maximum 50 characters")
	}

	// Check status length is not more than 50 characters
	if len(status) > 50 {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Status is too long, maximum 50 characters")
	}

	// Initialize sql.NullString for optional fields
	serialNumberNullStr := sql.NullString{}
	sizeNullStr := sql.NullString{}
	descriptionNullStr := sql.NullString{}
	statusNullStr := sql.NullString{}

	// convert serial number to sql.NullString if provided, else set to null
	if serialNumber == "" {
		serialNumberNullStr.String = serialNumber
		serialNumberNullStr.Valid = false
	} else {
		serialNumberNullStr.String = serialNumber
		serialNumberNullStr.Valid = true
	}

	// convert size to sql.NullString if provided, else set to null
	if size == "" {
		sizeNullStr.String = size
		sizeNullStr.Valid = false
	} else {
		sizeNullStr.String = size
		sizeNullStr.Valid = true
	}

	// convert description to sql.NullString if provided, else set to null
	if description == "" {
		descriptionNullStr.String = description
		descriptionNullStr.Valid = false
	} else {
		descriptionNullStr.String = description
		descriptionNullStr.Valid = true
	}

	// convert status to sql.NullString if provided, else set to null
	if status == "" {
		statusNullStr.String = status
		statusNullStr.Valid = false
	} else {
		statusNullStr.String = status
		statusNullStr.Valid = true
	}

	// Create new emergency device model
	emergencyDevice := &models.EmergencyDevice{
		RoomID:                roomID,
		EmergencyDeviceTypeID: emergencyDeviceTypeID,
		ExtinguisherTypeID:    extinguisherTypeID,
		SerialNumber:          serialNumberNullStr,
		ManufactureDate:       manufactureDate,
		LastInspectionDate:    lastInspectionDate,
		Size:                  sizeNullStr,
		Description:           descriptionNullStr,
		Status:                statusNullStr,
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
