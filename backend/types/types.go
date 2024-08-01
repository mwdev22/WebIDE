package types

import "github.com/go-playground/validator/v10"

type User struct {
	ID        uint64 `json:"id" validate:"required"`
	Username  string `json:"username" validate:"required"`
	GithubURL string `json:"git_url" validate:"required"`
}

type Repo struct {
	Name    string `json:"name" validate:"required"`
	Private bool   `json:"private" validate:"required"`
	UserID  int    `json:"user_id" validate:"required"`
}

type File struct {
	Name         string `json:"name" validate:"required"`
	Content      string `json:"content" validate:"required"`
	RepositoryID int    `json:"repo_id" validate:"required"`
}

var Validator = validator.New()
