package storage

import (
	"fmt"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Name         string
	Content      string
	RepositoryID int
}

type FileStore struct {
	db *gorm.DB
}

func (s *FileStore) GetFileByID(id int) (*File, error) {
	var file File
	if err := s.db.Where("ID = ?", id).First(&file).Error; err != nil {
		return nil, fmt.Errorf("failed to get file with id %v, %s", id, err)
	}
	return &file, nil
}
