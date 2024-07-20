package storage

import (
	"fmt"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name  string
	Bio   string
	Files File
}

type UserStore struct {
	db *gorm.DB
}

func (s *UserStore) GetUsers() ([]User, error) {
	users := s.db.Find(&User{}).Error
	fmt.Println(users)
	return []User{}, nil
}

func (s *UserStore) GetUserByID(id int) (*User, error) {
	var user User
	if err := s.db.Where("ID = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to get user with id %v, %s", id, err)
	}
	return &user, nil
}
