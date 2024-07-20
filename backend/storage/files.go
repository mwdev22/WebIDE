package storage

import (
	"gorm.io/gorm"
)

type File struct {
	gorm.Model

	name    string
	content string
	UserID  int
}

type FileStore struct {
	db *gorm.DB
}
