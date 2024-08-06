package storage

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type File struct {
	BaseModel
	Name         string `json:"name"`
	Content      string `json:"content"`
	RepositoryID int    `json:"repository_id"`

	// default fields display modifications
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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
