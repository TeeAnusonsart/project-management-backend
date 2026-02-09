package domain

type User struct {
	ID          uint
	Username    string
	Name        string
	Password    string
	Role        UserRole
	Email       string
	ProfilePath string
}

type UserRole string

const (
	RoleAdmin UserRole = "ADMIN"
	RoleUser  UserRole = "USER"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}