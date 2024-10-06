package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Sender    string `gorm:"not null" json:"sender"`
	Recipient string `gorm:"not null" json:"recipient"`
	Content   string `gorm:"not null" json:"content"`
	Timestamp string `gorm:"not null" json:"timestamp"`
}
