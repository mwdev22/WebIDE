package storage

import (
	"fmt"

	"github.com/mwdev22/WebIDE/cmd/types"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	ID           uint         `gorm:"primary_key" json:"id"`
	Username     string       `gorm:"not null" json:"username"`
	Bio          string       `json:"bio"`
	GithubURL    string       `gorm:"not null" json:"git_url"`
	Repositories []Repository `gorm:"foreignKey:OwnerID" json:"repositories"`
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

func (s *UserStore) GetUserByID(id uint) (*User, error) {
	var user User
	if err := s.db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user with id %v, %s", id, err)
	}
	return &user, nil
}

func (s *UserStore) CreateUser(data *types.UserPayload) error {
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
