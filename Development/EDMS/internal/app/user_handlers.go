package app

import (
	"database/sql"
	"log"
	"net/http"
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

	// Convert the user ID to an integer
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return a.handleError(c, http.StatusBadRequest, "Invalid user ID", err)
	}

	// check if updated username is unique
	existingUser, err := a.DB.GetUserByUsername(username)
	if err == nil {
		if existingUser.UserID != userIDInt {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Username already exists")
		}
	} else if err != sql.ErrNoRows {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching user", err)
	}

	// check if updated email is unique
	existingUser, err = a.DB.GetUserByEmail(email)
	if err == nil {
		if existingUser.UserID != userIDInt {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Email already exists")
		}
	} else if err != sql.ErrNoRows {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching user", err)
	}

	// Check if the user is trying to edit their own account
	if currentUserID == userID {
		// Check if the user is trying to change their role
		if defaultAdmin == "true" && role != "Admin" {
			return c.Redirect(http.StatusSeeOther, "/admin?error=Cannot change role on default admin account.")
		}

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
			// Check for errors
			// iF there is an error, return an error message
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

// func to HandleDeleteUser
func (a *App) HandleDeleteUser(c echo.Context) error {
	// Check if request if a POST request
	if c.Request().Method != http.MethodDelete {
		return c.JSON(http.StatusMethodNotAllowed, map[string]string{
			"error":       "Method not allowed",
			"redirectURL": "/admin?error=Method not allowed",
		})
	}

	// Get the user ID from the URL
	userID := c.Param("id")
	currentUserID := c.QueryParam("currentUserId")

	// convert the user ID & currentuserid to an integers
	userIDInt, err := strconv.Atoi(userID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid user ID",
			"redirectURL": "/admin?error=Invalid user ID" + userID,
		})
	}

	currentUserIDInt, err := strconv.Atoi(currentUserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid user ID",
			"redirectURL": "/admin?error=Invalid user ID" + currentUserID,
		})
	}

	// Get the user by ID
	user, err := a.DB.GetUserByID(userIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching user",
			"redirectURL": "/admin?error=Error fetching user",
		})
	}

	// Check if the user is trying to delete the default admin
	if user.DefaultAdmin {
		return c.JSON(http.StatusOK, map[string]string{
			"error":       "Cannot delete the default admin",
			"redirectURL": "/admin?error=Cannot delete the default admin",
		})
	}

	// Delete the user from the database
	err = a.DB.DeleteUser(userIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error deleting user",
			"redirectURL": "/admin?error=Error deleting user",
		})
	}

	// Check if the user is trying to delete their own account
	if currentUserIDInt == userIDInt {
		// Log the user out
		return c.JSON(http.StatusOK, map[string]string{
			"message":     "User deleted successfully",
			"redirectURL": "/logout?message=User deleted successfully. Please log in again",
		})
	}

	// Respond to the client
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "User deleted successfully",
		"redirectURL": "/admin?message=User deleted successfully",
	})
}
