package models

type User struct {
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	DefaultAdmin  bool   `json:"default_admin"`
	CurrentUserID int    `json:"current_user_id"`
}

type UserDto struct {
	CurrentUserID   string `json:"current_user_id"`
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	DefaultAdmin    string `json:"default_admin"`
}
