package models

import (
	"gorm.io/gorm"
)

// User represents the 'users' table in the database
type User struct {
<<<<<<< HEAD
	gorm.Model           // Contains ID, CreatedAt, UpdatedAt, DeletedAt
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Do not return password in JSON
	Email    string `gorm:"unique;not null" json:"email"`
	FullName string `json:"full_name"`
	Role     string `gorm:"default:customer" json:"role"` // User authorization roles
=======
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Không trả về password
	Email    string `gorm:"unique;not null" json:"email"`
	FullName string `json:"full_name"`
	Role     string `gorm:"default:customer" json:"role"` // Thêm phân quyền
>>>>>>> d37f6887b64410bea1d20c4a06ee4ef793a1113a
}
