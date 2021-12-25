package model

import "gorm.io/gorm"

/**
User model
利用gorm.Model
包含 ID, CreatedAt, UpdatedAt, DeletedAt
*/
type User struct {
	gorm.Model
	// `gorm:"type:varchar(20);not null"`
	Name     string
	Phone    string
	Password string
}
