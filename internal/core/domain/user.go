package domain

type User struct {
	Email       string
	Name string
	Role        UserRole
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

type UserRepository interface {
	Create(user *User) error
	FindAll() ([]*User, error)
	// FindByUID(uid uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Delete(email string) error

	// UpdatePassword(id uint, hashedPassword string) error
	// UpdateProfilePath(id uint, path string) error
}
