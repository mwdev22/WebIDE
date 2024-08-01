package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/backend/types"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint64 `gorm:"primarykey"`
	Username     string `gorm:"not null"`
	Bio          string
	GithubURL    string `gorm:"not null"`
	Repositories []Repository
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) GetAllUsers() ([]User, error) {
	var users []User
	if err := s.db.Find(users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %s", err)
	}

	fmt.Println(users)
	return users, nil
}

func (s *UserStore) GetUserByID(id uint64) (*User, error) {
	var user User
	if err := s.db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user with id %v, %s", id, err)
	}
	return &user, nil
}

func (s *UserStore) CreateUser(data *types.User) error {
	var newUser = User{
		ID:        data.ID,
		Username:  data.Username,
		GithubURL: data.GithubURL,
	}
	if err := s.db.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}
