package storage

import (
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	Name string

	UserID int
	Files  []File
}

type File struct {
	gorm.Model

	Name         string
	Content      string
	RepositoryID int
}

type FileStore struct {
	db *gorm.DB
}
