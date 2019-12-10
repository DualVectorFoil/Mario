package model

import "github.com/jinzhu/gorm"

type ProfileInfo struct {
	gorm.Model
	PhoneNum     string `gorm:"not null;index:idx_phone_num"`
	UserName     string `gorm:"not null"`
	PWDEncoded   string `gorm:"not null"`
	RegisteredAt int64  `gorm:"not null"`
	LastLoginAt  int64  `gorm:"not null"`
	Locale       string `gorm:"not null"`
}
