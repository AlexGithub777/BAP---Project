package app

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/config"
	"github.com/AlexGithub777/BAP---Project/Development/EDMS/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)

// CustomClaims represents JWT custom claims
type CustomClaims struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	DefaultAdmin bool   `json:"default_admin"`
	jwt.RegisteredClaims
}

// HandlePostForgotPassword handles the forgot password form submission
func (a *App) HandlePostForgotPassword(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "forgot_password.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	email := c.FormValue("email")

	// Check if the email exists in the database
	user, err := a.DB.GetUserByEmail(email)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/forgot-password?error=Email%20not%20found")
	}

	// Generate a new password
	newPassword, err := password.Generate(10, 8, 2, false, false)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/forgot-password?error=Could%20not%20generate%20password")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password
	if err := a.DB.UpdatePassword(user.UserID, string(hashedPassword)); err != nil {
		return c.Redirect(http.StatusSeeOther, "/forgot-password?error=Could%20not%20update%20password")
	}

	emailmessage := fmt.Sprintf("Could not send email but Your new password is: %s", newPassword)

	a.handleLogger(emailmessage)

	// Send the new password to the user's email
	if err := sendPasswordResetEmail(email, user.Username, newPassword); err != nil {
		return c.Redirect(http.StatusSeeOther, "/?message="+emailmessage)
	}
	message := fmt.Sprintf("Password reset successful. Check your %s for the new password.", email)

	// Render the login page with a success message
	return c.Redirect(http.StatusSeeOther, "/?message="+message)
}

// HandlePostRegister handles the register form submission
func (a *App) HandlePostRegister(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "register.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")
	confirmpassword := c.FormValue("confirm-password")

	// Validate the form data
	if username == "" || email == "" || password == "" || confirmpassword == "" {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "All fields are required",
		})
	}

	// Validate username
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{6,}$`)
	if !usernameRegex.MatchString(username) {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Username must be at least 6 characters and contain only letters, numbers, and underscores",
		})
	}

	// Validate email
	emailRegex := regexp.MustCompile(`[^@\s]+@[^@\s]+\.[^@\s]+`)
	if !emailRegex.MatchString(email) {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Invalid email address",
		})

	}

	// Validate password confirmation
	if password != confirmpassword {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Passwords do not match",
		})
	}

	// Validate password
	passwordLengthRegex := regexp.MustCompile(`.{8,}`)
	passwordDigitRegex := regexp.MustCompile(`[0-9]`)
	passwordSpecialCharRegex := regexp.MustCompile(`[!@#$%^&*]`)
	passwordCapitalLetterRegex := regexp.MustCompile(`[A-Z]`)

	if !passwordLengthRegex.MatchString(password) || !passwordDigitRegex.MatchString(password) || !passwordSpecialCharRegex.MatchString(password) || !passwordCapitalLetterRegex.MatchString(password) {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Password must contain at least one number, one special character, one capital letter, and be at least 8 characters long",
		})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Could not hash password",
		})
	}

	// Check if the user or email already exists
	if _, err := a.DB.GetUserByUsername(username); err == nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Username already exists",
		})
	}

	if _, err := a.DB.GetUserByEmail(email); err == nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Email already exists",
		})
	}

	// Create a new user
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := a.DB.CreateUser(&user); err != nil {
		return c.Render(http.StatusOK, "register.html", map[string]interface{}{
			"error": "Could not create user",
		})
	}

	// Generate a success message
	// Maybe send a welcome email here
	message := fmt.Sprintf("Registration successful. Please login with your username: %s", username)
	return c.Redirect(http.StatusSeeOther, "/?message="+message)
}

// HandleGetLogin serves the home page
func (a *App) HandleGetLogin(c echo.Context) error {
	// Check if request is a POST request
	if c.Request().Method != http.MethodGet {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	// Check if the user is already logged in
	cookie, err := c.Cookie("token")
	if err == nil && cookie.Value != "" {
		// Parse the JWT token
		token, err := parseToken(cookie.Value)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*CustomClaims); ok {
				// Put the claims data in the context
				c.Set("user", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Set("email", claims.Email)
				c.Set("default_admin", claims.DefaultAdmin)

				// User is already logged in, redirect to the dashboard
				return c.Redirect(http.StatusSeeOther, "/dashboard")
			}
		}
	}

	// If no valid token, render login page
	return c.Render(http.StatusOK, "index.html", nil)
}

// HandlePostLogin handles user login
func (a *App) HandlePostLogin(c echo.Context) error {
	// Check if request if a GET request
	if c.Request().Method != http.MethodPost {
		return c.Render(http.StatusMethodNotAllowed, "index.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	username := c.FormValue("username")
	password := c.FormValue("password")
	remember := c.FormValue("remember")

	// Validate the user's credentials
	user, err := a.DB.GetUserByUsername(username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"error": "Invalid username or password",
		})
	}

	// Determine expiration time based on "remember" checkbox
	expiresAt := time.Now().Add(72 * time.Hour) // Default expiration time is 3 days
	if remember == "on" {
		expiresAt = time.Now().Add(30 * 24 * time.Hour)
	}

	// Generate token
	token, err := GenerateToken(user, expiresAt)
	if err != nil {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"error": "Could not generate token",
		})
	}

	// Set the token as a cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	c.Set("user", token)

	// Return the token in the response as well
	return c.Redirect(http.StatusFound, "/dashboard")
}

// HandleGetLogout logs the user out
func (a *App) HandleGetLogout(c echo.Context) error {
	// Get message from query string
	message := c.QueryParam("message")
	// Check if request if a POST request
	if c.Request().Method != http.MethodGet {
		return c.Render(http.StatusMethodNotAllowed, "dashboard.html", map[string]interface{}{
			"error": "Method not allowed",
		})
	}

	// Clear JWT or session cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	c.SetCookie(cookie)

	// if message is empty, don't show it
	if message == "" {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	// Redirect the user to the login page
	return c.Redirect(http.StatusSeeOther, "/?message="+message)
}

// GenerateToken generates a JWT token
func GenerateToken(user *models.User, expiresAt time.Time) (string, error) {
	claims := &CustomClaims{
		UserID:       strconv.Itoa(user.UserID),
		Email:        user.Email,
		Username:     user.Username,
		Role:         user.Role,
		DefaultAdmin: user.DefaultAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.LoadConfig().JWTSecret
	return token.SignedString([]byte(secret))
}

// parseToken parses and validates the JWT token
func parseToken(tokenString string) (*jwt.Token, error) {
	secret := config.LoadConfig().JWTSecret
	return jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

// sendPasswordResetEmail sends a password reset email
func sendPasswordResetEmail(email, username, newPassword string) error {
	m := gomail.NewMessage()
	m.Reset()
	m.SetHeader("From", "MS_kAX6BH@trial-jpzkmgqev0vg059v.mlsender.net")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "EDMS PASSWORD RESET")
	m.SetBody("text/plain", "Your Username is "+username+", Your new password is: "+newPassword)
	m.AddAlternative("text/html", `<html><body style="font-family: Arial, sans-serif; padding: 20px;">
		<h2 style="color: #333;">EDMS PASSWORD RESET</h2>
		<p style="margin-top: 20px;">Your Username is <strong>`+username+`</strong></p>
		<p>Your new password is: <strong>`+newPassword+`</strong></p>
	</body></html>`)

	d := gomail.NewDialer("smtp.gmail.com", 587, "edms4758@gmail.com", "dksz kwoc zjiv cqub")

	fmt.Println("Sending email to", email)
	fmt.Println("Username:", username)
	fmt.Println("New Password:", newPassword)
	return d.DialAndSend(m)
}
