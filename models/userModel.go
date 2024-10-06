package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `gorm:"uniqueIndex;not null" json:"username"`
	Password    string `gorm:"not null" json:"-"`
	PhoneNumber string `gorm:"uniqueIndex;not null" json:"phone_number"`
}
