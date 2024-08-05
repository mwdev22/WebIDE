package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/backend/types"
	"gorm.io/gorm"
)

type Repository struct {
	gorm.Model
	ID      int    `gorm:"primarykey"`
	Name    string `json:"name"`
	Private bool   `json:"private"`
	UserID  uint   `json:"user_id"`
	Files   []File `json:"files"`
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

func (s *RepoStore) GetReposByUserID(id int) ([]*Repository, error) {
	var repos []*Repository

	// Fetch all repositories where UserID matches the given id
	if err := s.db.Where("user_id = ?", id).Find(&repos).Error; err != nil {
		return nil, fmt.Errorf("failed to get repositories for user_id %v: %s", id, err)
	}

	return repos, nil
}

func (s *RepoStore) CreateRepo(data types.Repo) (int, error) {
	var repo = Repository{
		Name:    data.Name,
		Private: data.Private,
		UserID:  data.UserID,
	}

	if err := s.db.Create(&repo).Error; err != nil {
		return 0, err
	}

	return repo.ID, nil
}
