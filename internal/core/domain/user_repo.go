package domain

type UserRepository interface {
	Create(user *User) error
	FindAll() ([]*User, error)
	FindByID(id uint) (*User, error)
	Delete(id uint) error

	UpdatePassword(id uint, hashedPassword string) error
	UpdateProfilePath(id uint, path string) error
}
