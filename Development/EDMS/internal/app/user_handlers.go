package app

import (
	"database/sql"
	"net/http"
	"regexp"
	"strconv"

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
	user, err := a.DB.GetUserByUsername(username)
	if err != nil {
		return a.handleError(c, http.StatusInternalServerError, "Error fetching data", err)
	}

	// Return the results as JSON
	return c.JSON(http.StatusOK, user)
}

// func to HandleEditUser
func (a *App) HandlePutUser(c echo.Context) error {
	// Check if request is a PUT request
	if c.Request().Method != http.MethodPut {
		return c.Redirect(http.StatusSeeOther, "/admin?error=Method not allowed")
	}

	// Parse the user ID from the URL parameter
	userID := c.Param("id")

	// Parse form data from the request body
	var user models.UserDto
	if err := c.Bind(&user); err != nil {
		a.handleLogger("Invalid request payload")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid request payload",
			"redirectURL": "/admin?error=Invalid request payload",
		})
	}

	user.UserID = userID

	// Validate input
	if user.Username == "" || user.Email == "" || user.Role == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Username, email, and role are required",
			"redirectURL": "/admin?error=Username, email, and role are required",
		})
	}

	// Validate username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
	if !usernameRegex.MatchString(user.Username) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Username must be at least 6 characters long and contain only letters, numbers, and underscores",
			"redirectURL": "/admin?error=Username must be at least 6 characters long and contain only letters, numbers, and underscores",
		})
	}

	// Validate email
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(user.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid email address",
			"redirectURL": "/admin?error=Invalid email address",
		})
	}

	// Validate role
	if user.Role != "Admin" && user.Role != "User" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid role",
			"redirectURL": "/admin?error=Invalid role",
		})
	}

	// Convert the user ID to an integer
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		a.handleLogger("Invalid user ID")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":       "Invalid user ID",
			"redirectURL": "/admin?error=Invalid user ID",
		})

	}

	// check if updated user.Username is unique
	existingUser, err := a.DB.GetUserByUsername(user.Username)
	if err == nil {
		if existingUser.UserID != userIDInt {
			return c.JSON(http.StatusOK, map[string]string{
				"error":       "Username already exists",
				"redirectURL": "/admin?error=Username already exists",
			})
		}
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching user",
			"redirectURL": "/admin?error=Error fetching user",
		})
	}

	// check if updated email is unique
	existingUser, err = a.DB.GetUserByEmail(user.Email)
	if err == nil {
		if existingUser.UserID != userIDInt {
			return c.JSON(http.StatusOK, map[string]string{
				"error":       "Email already exists",
				"redirectURL": "/admin?error=Email already exists",
			})
		}
	} else if err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":       "Error fetching user",
			"redirectURL": "/admin?error=Error fetching user",
		})
	}

	// Check if the user is trying to edit their own account
	if user.CurrentUserID == userID {
		// Check if the user is trying to change their role
		if user.DefaultAdmin == "true" && user.Role != "Admin" {
			return c.JSON(http.StatusOK, map[string]string{
				"error": "Cannot change role of default admin",
			})
		}

		// parse tthe password
		password := user.Password
		confirmedPassword := user.ConfirmPassword

		if password == "" {
			// update the user without changing the password
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error":       "Invalid user ID",
					"redirectURL": "/admin?error=Invalid user ID",
				})
			}

			user := &models.User{
				UserID:   userIDInt,
				Username: user.Username,
				Email:    user.Email,
				Role:     user.Role,
			}

			// Update the user in the database
			err = a.DB.UpdateUser(user)
			// Check for errors
			// iF there is an error, return an error message
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":       "Error updating user",
					"redirectURL": "/admin?error=Error updating user",
				})
			}

			// Log the user out
			return c.JSON(http.StatusOK, map[string]string{
				"message":     "User details updated successfully. Please log in again",
				"redirectURL": "/logout?message=User details updated successfully. Please log in again",
			})
		} else {
			// Validate password
			if password != confirmedPassword {
				return c.JSON(http.StatusOK, map[string]string{
					"error":       "Passwords do not match",
					"redirectURL": "/admin?error=Passwords do not match",
				})
			}

			passwordLengthRegex := regexp.MustCompile(`.{8,}`)
			passwordDigitRegex := regexp.MustCompile(`[0-9]`)
			passwordSpecialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)
			passwordCapitalLetterRegex := regexp.MustCompile(`[A-Z]`)

			if !passwordLengthRegex.MatchString(password) || !passwordDigitRegex.MatchString(password) || !passwordSpecialCharRegex.MatchString(password) || !passwordCapitalLetterRegex.MatchString(password) {
				return c.JSON(http.StatusOK, map[string]string{
					"error":       "Password must be at least 8 characters long and contain at least one digit, special character, and capital letter",
					"redirectURL": "/admin?error=Password must be at least 8 characters long and contain at least one digit, special character, and capital letter",
				})
			}

			// Hash the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":       "Error hashing password",
					"redirectURL": "/admin?error=Error hashing password",
				})
			}

			// make a user model
			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error":       "Invalid user ID",
					"redirectURL": "/admin?error=Invalid user ID",
				})
			}

			user := &models.User{
				UserID:   userIDInt,
				Username: user.Username,
				Email:    user.Email,
				Password: string(hashedPassword),
				Role:     user.Role,
			}

			// Update the user in the database
			err = a.DB.UpdateUserWithPassword(user)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error":       "Error updating user",
					"redirectURL": "/admin?error=Error updating user",
				})

			}
			// Log the user out
			return c.JSON(http.StatusOK, map[string]string{
				"message":     "User details updated successfully. Please log in again",
				"redirectURL": "/logout?message=User details updated successfully. Please log in again",
			})
		}

	} else {
		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error":       "Invalid user ID",
				"redirectURL": "/admin?error=Invalid user ID",
			})
		}

		// Check if tyring to update the default admin
		if user.DefaultAdmin == "true" {
			return c.JSON(http.StatusOK, map[string]string{
				"error":       "Cannot change role of default admin",
				"redirectURL": "/admin?error=Cannot change role of default admin",
			})
		}

		user := &models.User{
			UserID:   userIDInt,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		}

		// Update the user in the database
		err = a.DB.UpdateUser(user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error":       "Error updating user",
				"redirectURL": "/admin?error=Error updating user",
			})
		}
	}

	// Redirect to the admin page with a success message
	return c.JSON(http.StatusOK, map[string]string{
		"message":     "User details updated successfully",
		"redirectURL": "/admin?message=User details updated successfully",
	})
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
