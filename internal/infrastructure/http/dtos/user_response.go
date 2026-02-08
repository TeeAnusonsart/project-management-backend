package dtos

type UserResponse struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ProfilePath string `json:"profile_path"`
	Role        string `json:"role"`
}
