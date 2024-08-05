package types

import "github.com/go-playground/validator/v10"

type User struct {
	ID        uint   `json:"id" validate:"required"`
	Username  string `json:"username" validate:"required"`
	GithubURL string `json:"git_url" validate:"required"`
}

type Repo struct {
	Name    string `json:"name" validate:"required"`
	Private bool   `json:"private" validate:"required"`
	UserID  uint   `json:"user_id" validate:"required"`
}

type File struct {
	Name         string `json:"name" validate:"required"`
	Content      string `json:"content" validate:"required"`
	RepositoryID int    `json:"repo_id" validate:"required"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
