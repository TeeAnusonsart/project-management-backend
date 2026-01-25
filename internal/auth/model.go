package auth

import (
	"gorm.io/gorm"
)

type UserAccount struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"password"`
}
