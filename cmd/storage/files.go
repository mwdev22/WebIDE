package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/cmd/types"
	"gorm.io/gorm"
)

type File struct {
	BaseModel
	Name         string `json:"name"`
	Content      string `json:"content"`
	Extension    string `json:"extension"`
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

func (s *FileStore) CreateFile(data types.FilePayload) (uint, error) {
	var file = File{
		Name:         data.Name,
		Content:      data.Content,
		RepositoryID: data.RepositoryID,
		Extension:    data.Extension,
	}

	if err := s.db.Create(&file).Error; err != nil {
		return 0, err
	}

	return file.ID, nil
}

func (s *FileStore) UpdateFile(file *File) error {
	return s.db.Save(file).Error
}

func (s *FileStore) GetFilesByRepoID(id uint) ([]File, error) {
	var files []File
	if err := s.db.Where("repository_id = ?", id).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("failed to get files for repository id %v, %s", id, err)
	}
	return files, nil
}
