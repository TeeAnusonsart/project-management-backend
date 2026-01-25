package auth

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) getByUsername(username string) (*UserAccount, error) {
	var user UserAccount
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *Repository) createUser(user *UserAccount) error {
	return r.db.Create(user).Error
}
