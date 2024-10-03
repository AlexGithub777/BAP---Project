package models

// UserT represents the users in the system
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
	CurrentUserID   string `json:"currentUserID"`
	UserID          string `json:"editUserID"`
	Username        string `json:"editUserUsername"`
	Email           string `json:"editUserEmail"`
	Role            string `json:"editUserRole"`
	Password        string `json:"editUserPassword"`
	ConfirmPassword string `json:"editUserConfirmPassword"`
	DefaultAdmin    string `json:"defaultAdmin"`
}
