package storage

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID      uint
	Name    string
	Profile Profile
}

type Profile struct {
	gorm.Model

	ID     uint
	UserID uint
	Bio    string
}
