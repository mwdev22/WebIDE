package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/backend/types"
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	Name    string
	Private bool
	UserID  int
	Files   []File
}

type RepoStore struct {
	db *gorm.DB
}

func NewRepoStore(db *gorm.DB) *RepoStore {
	return &RepoStore{
		db: db,
	}
}

func (s *RepoStore) GetRepoByID(id int) (*Repository, error) {
	var repo Repository
	if err := s.db.Where("ID = ?", id).First(&repo).Error; err != nil {
		return nil, fmt.Errorf("failed to get file with id %v, %s", id, err)
	}
	return &repo, nil
}

func (s *RepoStore) CreateRepo(data types.Repo) error {
	var repo = Repository{
		Name:    data.Name,
		Private: data.Private,
		UserID:  data.UserID,
	}
	if err := s.db.Create(&repo).Error; err != nil {
		return err
	}
	return nil
}
