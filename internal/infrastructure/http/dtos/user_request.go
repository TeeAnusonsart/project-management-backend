package dtos

type CreateUserRequest struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	ProfilePath string `json:"profile_path"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
