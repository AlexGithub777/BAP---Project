package app

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// HandleGetAllUsers fetches all users from the database and returns the results as JSON
func (a *App) HandleGetAllUsers(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	users, err := a.DB.GetAllUsers()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, users)
}

// HandleGetUserByUsername
func (a *App) HandleGetUserByUsername(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")

	}

	username := c.Param("username")
	log.Printf("Username: %s", username)
	user, err := a.DB.GetUserByUsername(username)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, user)
}

// func to HandleEditUser
func (a *App) HandleEditUser(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form
	err := c.Request().ParseForm()
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error parsing form")
	}

	// Get the form values
	currentUserID := c.FormValue("currentUserID")
	userID := c.FormValue("editUserID")
	username := strings.TrimSpace(c.FormValue("editUserUsername"))
	email := strings.TrimSpace(c.FormValue("editUserEmail"))
	role := strings.TrimSpace(c.FormValue("editUserRole"))
	defaultAdmin := c.FormValue("defaultAdmin")

	log.Printf("currentUserID: %s, userID: %s, username: %s, email: %s, role: %s", currentUserID, userID, username, email, role)

	// Validate input
	if username == "" || email == "" || role == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=All fields are required")
	}

	// Validate username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
	if !usernameRegex.MatchString(username) {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Username must be at least 6 characters long and contain only letters, numbers, and underscores")
	}

	// Validate email
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email) {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid email address")
	}

	// Validate role
	if role != "Admin" && role != "User" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid role")
	}

	// Check if the user is trying to edit their own account
	if currentUserID == userID {
		// parse tthe password
		password := c.FormValue("editUserPassword")
		confirmedPassword := c.FormValue("editUserConfirmPassword")

		if password == "" {
			// update the user without changing the password
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				return a.handleError(c, http.StatusBadRequest, "Invalid user ID", err)
			}

			user := &models.User{
				UserID:   userIDInt,
				Username: username,
				Email:    email,
				Role:     role,
			}

			// Update the user in the database
			err = a.DB.UpdateUser(user)
			if err != nil {
				return a.handleError(c, http.StatusInternalServerError, "Error updating user", err)
			}

			// Log the user out
			return c.Redirect(http.StatusSeeOther, "/logout?message=User details updated successfully. Please log in again")
		} else {
			// Validate password
			if password != confirmedPassword {
				return c.Redirect(http.StatusSeeOther, "/admin?error=Passwords do not match")
			}

			passwordLengthRegex := regexp.MustCompile(`.{8,}`)
			passwordDigitRegex := regexp.MustCompile(`[0-9]`)
			passwordSpecialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)
			passwordCapitalLetterRegex := regexp.MustCompile(`[A-Z]`)

			if !passwordLengthRegex.MatchString(password) || !passwordDigitRegex.MatchString(password) || !passwordSpecialCharRegex.MatchString(password) || !passwordCapitalLetterRegex.MatchString(password) {
				return c.Redirect(http.StatusSeeOther, "/admin?error=Password must be at least 8 characters long and contain at least one digit, one special character, and one capital letter")
			}

			// Hash the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			if err != nil {
				return a.handleError(c, http.StatusInternalServerError, "Error hashing password", err)
			}

			// make a user model
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				return a.handleError(c, http.StatusBadRequest, "Invalid user ID", err)
			}

			user := &models.User{
				UserID:   userIDInt,
				Username: username,
				Email:    email,
				Password: string(hashedPassword),
				Role:     role,
			}

			// Update the user in the database
			err = a.DB.UpdateUserWithPassword(user)
			if err != nil {
				return a.handleError(c, http.StatusInternalServerError, "Error updating user", err)

			}
			// Log the user out
			return c.Redirect(http.StatusSeeOther, "/logout?message=User details updated successfully. Please log in again")
		}

	} else {
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			return a.handleError(c, http.StatusBadRequest, "Invalid user ID", err)
		}

		// Check if tyring to update the default admin
		if defaultAdmin == "true" {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Cannot update the default admin")
		}

		user := &models.User{
			UserID:   userIDInt,
			Username: username,
			Email:    email,
			Role:     role,
		}

		// Update the user in the database
		err = a.DB.UpdateUser(user)
		if err != nil {
			return a.handleError(c, http.StatusInternalServerError, "Error updating user", err)
		}
	}

	// Redirect to the admin page with a success message
	return c.Redirect(http.StatusFound, "/admin?message=User updated successfully")
}

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

func (a *App) HandleGetAllDeviceTypes(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	emergencyDeviceTypes, err := a.DB.GetAllDeviceTypes()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, emergencyDeviceTypes)
}

func (a *App) HandleGetAllExtinguisherTypes(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	extinguisherTypes, err := a.DB.GetAllExtinguisherTypes()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, extinguisherTypes)
}

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

func (a *App) HandleGetAllSites(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	sites, err := a.DB.GetAllSites()
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, sites)
}

func (a *App) HamdleGetSiteByID(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Redirect(http.StatusSeeOther, "/dashboard?error=Method not allowed")
	}

	id := c.Param("id")
	site, err := a.DB.GetSiteByID(id)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, site)
}

func (a *App) HandlePostSite(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form, limiting upload size to 10MB
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error parsing form")
	}

	// Get the form values
	siteName := strings.TrimSpace(c.FormValue("addSiteName"))       // Trim whitespace
	siteAddress := strings.TrimSpace(c.FormValue("addSiteAddress")) // Trim whitespace

	// Validate input
	if siteName == "" || siteAddress == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=All fields are required")
	}

	// Validate site name & address length (site name should be less than 100 characters) (address should be less than 255 characters)
	if len(siteName) > 100 || len(siteAddress) > 255 {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name should be less than 100 characters and address should be less than 255 characters")
	}

	// Additional validation for siteName (allow only alphanumeric, spaces, hyphens, and underscores)
	if !regexp.MustCompile(`^[a-zA-Z0-9\s_-]+$`).MatchString(siteName) {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name can only contain letters, numbers, spaces, hyphens, and underscores")
	}

	// Check if site name is unique
	_, err = a.DB.GetSiteByName(siteName)
	if err == nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name already exists")
	} else if err != sql.ErrNoRows { // If the error is not sql.ErrNoRows, it's a database error
		return c.Redirect(http.StatusSeeOther, "/admin?error=Database error")
	}

	// Initialize filePath as an empty sql.NullString
	filePath := sql.NullString{String: "", Valid: false}

	// Retrieve the file from the form
	file, header, err := c.Request().FormFile("siteMapImgInput")
	if err == nil {
		defer file.Close()
		// Validate file extension
		// Create unique file name based on the site name
		fileExt := filepath.Ext(header.Filename)
		allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".svg": true}
		if !allowedExtensions[fileExt] {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid file type. Allowed types: jpg, jpeg, png, gif, svg")
		}
		// Define static directory for site maps
		staticDir := "./static/site_maps"
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			os.MkdirAll(staticDir, os.ModePerm) // Create directory if it doesn't exist
		}

		sanitizedSiteName := strings.ReplaceAll(siteName, " ", "_")
		sanitizedSiteName = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(sanitizedSiteName, "")
		fileName := filepath.Join(staticDir, sanitizedSiteName+fileExt)

		out, err := os.Create(fileName)
		if err != nil {
			return c.Redirect(http.StatusInternalServerError, "/admin?error=Error creating file")
		}
		defer out.Close()

		// Copy the uploaded file data to the new file
		_, err = io.Copy(out, file)
		if err != nil {
			return c.Redirect(http.StatusInternalServerError, "/admin?error=Error copying file")
		}

		// Save the relative path as a sql.NullString
		filePath = sql.NullString{String: "/static/site_maps/" + sanitizedSiteName + fileExt, Valid: true}
	}

	// Save site information and file path in the database
	site := &models.Site{
		SiteName:         siteName,
		SiteAddress:      siteAddress,
		SiteMapImagePath: filePath,
	}

	err = a.DB.AddSite(site)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error saving site", err)
	}

	// Redirect to the admin page with a success message
	return c.Redirect(http.StatusFound, "/admin?message=Site added successfully")
}

func (a *App) HandleEditSite(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the form, limiting upload size to 10MB
	err := c.Request().ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Error parsing form")
	}

	// Get the form values
	siteID := c.FormValue("editSiteID")
	siteName := strings.TrimSpace(c.FormValue("editSiteName"))       // Trim whitespace
	siteAddress := strings.TrimSpace(c.FormValue("editSiteAddress")) // Trim whitespace

	// Validate input
	if siteName == "" || siteAddress == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?error=All fields are required")
	}

	// Validate site name & address length (site name should be less than 100 characters) (address should be less than 255 characters)
	if len(siteName) > 100 || len(siteAddress) > 255 {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name should be less than 100 characters and address should be less than 255 characters")
	}

	// Additional validation for siteName (allow only alphanumeric, spaces, hyphens, and underscores)
	if !regexp.MustCompile(`^[a-zA-Z0-9\s_-]+$`).MatchString(siteName) {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name can only contain letters, numbers, spaces, hyphens, and underscores")
	}

	// Get the existing site by ID
	existingSite, err := a.DB.GetSiteByID(siteID)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching site", err)
	}

	// Check if updated site name is unique
	siteWithSameName, err := a.DB.GetSiteByName(siteName)
	if err != nil {
		if err != sql.ErrNoRows { // If the error is not sql.ErrNoRows, it's a database error
			return a.handleError(c, http.StatusInternalServerError, "Database error", err)
		}
	} else if siteWithSameName.SiteID != existingSite.SiteID {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Site name already exists")
	}

	// Initialize filePath as an empty sql.NullString
	filePath := sql.NullString{String: "", Valid: false}

	// Define static directory for site maps
	staticDir := "./static/site_maps"

	// Create sanitized site name
	sanitizedSiteName := strings.ReplaceAll(siteName, " ", "_")
	sanitizedSiteName = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(sanitizedSiteName, "")

	// Initialize file extension
	var fileExt string

	// Retrieve the file from the form
	file, header, err := c.Request().FormFile("siteMapImgInput")
	if err == nil {
		defer file.Close()

		// If a new image is being uploaded and the existing site has an image, delete the old image
		if existingSite != nil && existingSite.SiteMapImagePath.Valid {
			oldImagePath := "." + existingSite.SiteMapImagePath.String
			if err := os.Remove(oldImagePath); err != nil {
				return c.Redirect(http.StatusInternalServerError, "/admin?error=Error deleting old image")
			}
		}

		// Validate file extension
		// Create unique file name based on the site name
		fileExt = filepath.Ext(header.Filename)
		allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".svg": true}
		if !allowedExtensions[fileExt] {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Invalid file type. Allowed types: jpg, jpeg, png, gif, svg")
		}

		// Define static directory for site maps
		if _, err := os.Stat(staticDir); os.IsNotExist(err) {
			os.MkdirAll(staticDir, os.ModePerm) // Create directory if it doesn't exist
		}

		sanitizedSiteName = strings.ReplaceAll(siteName, " ", "_")
		sanitizedSiteName = regexp.MustCompile(`[^a-zA-Z0-9_-]`).ReplaceAllString(sanitizedSiteName, "")
		fileName := filepath.Join(staticDir, sanitizedSiteName+fileExt)

		out, err := os.Create(fileName)
		if err != nil {
			return c.Redirect(http.StatusInternalServerError, "/admin?error=Error creating file")
		}
		defer out.Close()

		// Copy the uploaded file data to the new file
		_, err = io.Copy(out, file)
		if err != nil {
			return c.Redirect(http.StatusInternalServerError, "/admin?error=Error copying file")
		}

		// Save the relative path as a sql.NullString
		filePath = sql.NullString{String: "/static/site_maps/" + sanitizedSiteName + fileExt, Valid: true}
	} else {
		fileExt = filepath.Ext(existingSite.SiteMapImagePath.String)
	}

	// check if the site name has changed and a map exists
	if siteName != existingSite.SiteName && existingSite.SiteMapImagePath.Valid {
		// Only rename the file if no new file is uploaded
		if err != nil {
			oldImagePath := "." + existingSite.SiteMapImagePath.String
			newImagePath := filepath.Join(staticDir, sanitizedSiteName+fileExt)
			if err := os.Rename(oldImagePath, newImagePath); err != nil {
				return c.Redirect(http.StatusInternalServerError, "/admin?error=Error renaming image")
			}

			// Update the file path
			filePath = sql.NullString{String: "/static/site_maps/" + sanitizedSiteName + fileExt, Valid: true}
		}
	}

	// Save site information and file path in the database
	var siteMapImagePath sql.NullString
	if filePath.Valid {
		siteMapImagePath = filePath
	} else {
		siteMapImagePath = existingSite.SiteMapImagePath
	}

	site := &models.Site{
		SiteID:           existingSite.SiteID,
		SiteName:         siteName,
		SiteAddress:      siteAddress,
		SiteMapImagePath: siteMapImagePath, // Use the existing path if no new file was uploaded
	}

	err = a.DB.UpdateSite(site)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error saving site", err)
	}

	// Respond to the client
	return c.Redirect(http.StatusFound, "/admin?message=Site updated successfully")
}

func (a *App) HandleDeleteSite(c echo.Context) error {
	// Check if request is not a delete request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/admin?error=Method not allowed",
		})
	}

	// Get the site ID from the URL
	siteID := c.Param("id")

	// Get the site by ID
	site, err := a.DB.GetSiteByID(siteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching site",
			"redirectURL": "/admin?error=Error fetching site",
		})
	}

	// Check if the site has any emergency devices
	emergencyDevices, err := a.DB.GetAllDevices(siteID, "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching emergency devices",
			"redirectURL": "/admin?error=Error fetching emergency devices",
		})
	}

	// Check if the site has any emergency devices
	if len(emergencyDevices) > 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot delete site with associated emergency devices",
			"redirectURL": "/admin?error=Cannot delete site with associated emergency devices",
		})
	}

	// Check if the site has any rooms
	rooms, err := a.DB.GetRoomsBySiteID(siteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching rooms",
			"redirectURL": "/admin?error=Error fetching rooms",
		})
	}

	// Check if the site has any rooms
	if len(rooms) > 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot delete site with associated rooms",
			"redirectURL": "/admin?error=Cannot delete site with associated rooms",
		})
	}

	// Handle foreign key constraints
	// Check if the site has any buildings
	buildings, err := a.DB.GetAllBuildings(siteID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching buildings",
			"redirectURL": "/admin?error=Error fetching buildings",
		})
	}

	if len(buildings) > 0 {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot delete site with associated buildings",
			"redirectURL": "/admin?error=Cannot delete site with associated buildings",
		})
	}

	// Delete the site from the database
	err = a.DB.DeleteSite(siteID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting site",
			"redirectURL": "/admin?error=Error deleting site",
		})
	}

	// Check if the site has a map image
	if site.SiteMapImagePath.Valid {
		// Delete the map image file
		imagePath := "." + site.SiteMapImagePath.String
		if err := os.Remove(imagePath); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":       "Error deleting site map image",
				"redirectURL": "/admin?error=Error deleting site map image",
			})

		}
	}

	// Respond to the client
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "Site deleted successfully",
		"redirectURL": "/admin?message=Site deleted successfully",
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
