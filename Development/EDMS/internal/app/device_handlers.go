package app

import (
	"log"
	"net/http"

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
	emergencyDeviceTypeIDStr := c.FormValue("emergency_device_type_id")
	extinguisherTypeIDStr := c.FormValue("extinguisher_type_id")
	serialNumber := c.FormValue("serial_number")
	manufactureDateStr := c.FormValue("manufacture_date")
	lastInspectionDateStr := c.FormValue("last_inspection_date")
	size := c.FormValue("size")
	description := c.FormValue("description")
	status := c.FormValue("status")
	a.handleLogger("Room: " + roomIDStr)
	a.handleLogger("extinguisher_type_id: " + extinguisherTypeIDStr)
	a.handleLogger("serial_number: " + serialNumber)
	a.handleLogger("manufacture_date: " + manufactureDateStr)
	a.handleLogger("last_inspection_date: " + lastInspectionDateStr)
	a.handleLogger("size: " + size)
	a.handleLogger("description: " + description)
	a.handleLogger("status: " + status)
	a.handleLogger("extinguisher type: " + extinguisherTypeIDStr)
	//Validating device details
	if roomIDStr == "" || emergencyDeviceTypeIDStr == "" {
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}
	/*
		//validate room_id
		_, err = a.DB.GetRoomByID(roomIDStr)
		if err == nil {
			return c.Redirect(http.StatusSeeOther, "/dashboard")
		} else if err != sql.ErrNoRows { // If the error is not sql.ErrNoRows, it's a database error
			return c.Render(http.StatusInternalServerError, "admin.html", map[string]interface{}{
				"error": "Database error",
			})
		}


		/* Save site information and file path in the database
		site := &models.Site{
			SiteName:         siteName,
			SiteAddress:      siteAddress,
			SiteMapImagePath: filePath,
		}

		err = a.DB.AddDevice(site)
		if err != nil {
			return a.handleError(c, http.StatusInternalServerError, "Error saving site", err)
		}
	*/
	return c.Render(http.StatusOK, "dashboard.html", map[string]interface{}{
		"message": "Device added successfully",
	})
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
