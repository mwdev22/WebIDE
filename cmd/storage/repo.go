package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/cmd/types"
	"gorm.io/gorm"
)

type Repository struct {
	BaseModel
	Name    string `json:"name"`
	Private bool   `json:"private"`
	Readme  string `json:"readme"`
	OwnerID uint   `json:"owner_id"`
	Owner   User   `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
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
	if err := s.db.Preload("Files").First(&repo, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get file with id %v, %s", id, err)
	}
	return &repo, nil
}

func (s *RepoStore) GetReposByUserID(id int) ([]*Repository, error) {
	var repos []*Repository

	// Fetch all repositories where UserID matches the given id
	if err := s.db.Where("owner_id = ?", id).Find(&repos).Error; err != nil {
		return nil, fmt.Errorf("failed to get repositories for user_id %v: %s", id, err)
	}

	return repos, nil
}

func (s *RepoStore) CreateRepo(data types.RepoPayload) (uint, error) {
	var repo = Repository{
		Name:    data.Name,
		Private: data.Private,
		OwnerID: data.UserID,
	}

	if err := s.db.Create(&repo).Error; err != nil {
		return 0, err
	}

	return repo.ID, nil
}

func (s *RepoStore) UpdateRepo(repo *Repository) error {
	return s.db.Save(repo).Error
}
