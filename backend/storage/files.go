package storage

import (
	"fmt"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	ID           int    `gorm:"primarykey" json:"id"`
	Name         string `json:"name"`
	Content      string `json:"content"`
	RepositoryID int    `json:"repository_id"`
}

type FileStore struct {
	db *gorm.DB
}

func NewFileStore(db *gorm.DB) *FileStore {
	return &FileStore{
		db: db,
	}
}

func (s *FileStore) GetFileByID(id int) (*File, error) {
	var file File
	if err := s.db.Where("ID = ?", id).First(&file).Error; err != nil {
		return nil, fmt.Errorf("failed to get file with id %v, %s", id, err)
	}
	return &file, nil
}
