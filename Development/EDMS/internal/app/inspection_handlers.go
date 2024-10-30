package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/labstack/echo/v4"
)

// HandleGetInspectionsByDeviceID fetches all inspections by device ID
func (a *App) HandleGetAllInspectionsByDeviceID(c echo.Context) error {
	// Check if request is a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

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
	// Check if request is a GET request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

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

// Helper function to parse "on" as true and anything else (including "") as false
func parseCheckbox(value string) bool {
	return value == "on"
}

// HandlePostInspection creates a new inspection
func (a *App) HandlePostInspection(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	inspection := new(models.Inspection)

	// Parse required fields first
	inspectionDate := c.FormValue("inspection_date")
	notes := c.FormValue("notes")
	deviceID, err := strconv.Atoi(c.FormValue("device_id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid Device ID")
	}

	userId, err := strconv.Atoi(c.FormValue("user_id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid User ID")
	}
	// Usage example for each form field with checkboxes
	inspection.IsConspicuous.Bool = parseCheckbox(c.FormValue("isConspicuous"))
	inspection.IsAccessible.Bool = parseCheckbox(c.FormValue("isAccessible"))
	inspection.IsAssignedLocation.Bool = parseCheckbox(c.FormValue("isAssignedLocation"))
	inspection.IsSignVisible.Bool = parseCheckbox(c.FormValue("isSignVisible"))
	inspection.IsAntiTamperDeviceIntact.Bool = parseCheckbox(c.FormValue("isAntiTamperDeviceIntact"))
	inspection.IsSupportBracketSecure.Bool = parseCheckbox(c.FormValue("isSupportBracketSecure"))
	inspection.WorkOrderRequired.Bool = parseCheckbox(c.FormValue("workOrderRequired"))
	inspection.AreOperatingInstructionsClear.Bool = parseCheckbox(c.FormValue("areOperatingInstructionsClear"))
	inspection.IsMaintenanceTagAttached.Bool = parseCheckbox(c.FormValue("isMaintenanceTagAttached"))
	inspection.IsNoExternalDamage.Bool = parseCheckbox(c.FormValue("isNoExternalDamage"))
	inspection.IsChargeGaugeNormal.Bool = parseCheckbox(c.FormValue("isChargeGaugeNormal"))
	inspection.IsReplaced.Bool = parseCheckbox(c.FormValue("isReplaced"))
	inspection.AreMaintenanceRecordsComplete.Bool = parseCheckbox(c.FormValue("areMaintenanceRecordsComplete"))
	inspection_status := c.FormValue("inspection_status")

	// Log the inspection details
	fmt.Println("Inspection Date:", inspectionDate)
	fmt.Println("Notes:", notes)
	fmt.Println("Device ID:", deviceID)
	fmt.Println("User ID:", userId)
	fmt.Println("IsConspicuous:", inspection.IsConspicuous.Bool)
	fmt.Println("IsAccessible:", inspection.IsAccessible.Bool)
	fmt.Println("IsAssignedLocation:", inspection.IsAssignedLocation.Bool)
	fmt.Println("IsSignVisible:", inspection.IsSignVisible.Bool)
	fmt.Println("IsAntiTamperDeviceIntact:", inspection.IsAntiTamperDeviceIntact.Bool)
	fmt.Println("IsSupportBracketSecure:", inspection.IsSupportBracketSecure.Bool)
	fmt.Println("WorkOrderRequired:", inspection.WorkOrderRequired.Bool)
	fmt.Println("AreOperatingInstructionsClear:", inspection.AreOperatingInstructionsClear.Bool)
	fmt.Println("IsMaintenanceTagAttached:", inspection.IsMaintenanceTagAttached.Bool)
	fmt.Println("IsNoExternalDamage:", inspection.IsNoExternalDamage.Bool)
	fmt.Println("IsChargeGaugeNormal:", inspection.IsChargeGaugeNormal.Bool)
	fmt.Println("IsReplaced:", inspection.IsReplaced.Bool)
	fmt.Println("AreMaintenanceRecordsComplete:", inspection.AreMaintenanceRecordsComplete.Bool)
	fmt.Println("InspectionStatus:", inspection_status)

	// Validate required fields
	if inspectionDate == "" || deviceID == 0 || userId == 0 || inspection_status == "" {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid request payload")
	}

	// Check if the device ID exists
	_, err = a.DB.GetDeviceByID(deviceID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid Device ID")
	}

	// Check if the user ID exists
	_, err = a.DB.GetUserByID(userId)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid User ID")
	}

	// Parse and validate inspection date
	formattedInspectionDate, err := time.Parse("2006-01-02", inspectionDate)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Invalid Inspection Date")
	}

	// Create the sql.NullTime struct
	nullTimeDate := sql.NullTime{
		Time:  formattedInspectionDate,
		Valid: true,
	}

	// Check if date is in the future
	if time.Now().Before(formattedInspectionDate) {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Inspection Date cannot be in the future")
	}

	// Validate notes length is less than 255 characters
	if len(notes) > 255 {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Notes must be less than 255 characters")
	}

	// Set the remaining inspection fields
	inspection.InspectionDate = nullTimeDate
	inspection.Notes.String = notes
	inspection.EmergencyDeviceID = deviceID
	inspection.UserID = userId
	inspection.InspectionStatus = inspection_status

	// Log the inspection details (consider using structured logging)
	a.handleLogger(fmt.Sprintf("New inspection submission: deviceID=%d, userID=%d, date=%s",
		deviceID, userId, inspectionDate))

	// Add the inspection to the database
	err = a.DB.AddInspection(inspection)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/dashboard?message=Inspection added successfully")
}
