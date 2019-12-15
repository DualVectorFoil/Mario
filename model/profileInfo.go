package model

import "github.com/jinzhu/gorm"

type ProfileInfo struct {
	gorm.Model
	PhoneNum     string `gorm:"not null;index:idx_phone_num"`
	AvatarUrl    string
	UserName     string `gorm:"not null"`
	Locale       string
	Bio          string
	Followers    int32
	Following    int32
	ArtworkCount int32
	PWD          string `gorm:"not null"`
	RegisteredAt int64  `gorm:"not null"`
	LastLoginAt  int64  `gorm:"not null"`
}
