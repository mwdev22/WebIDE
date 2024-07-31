package types

import "github.com/go-playground/validator/v10"

type User struct {
	GithubID  uint64
	Username  string
	GithubURL string
}

var Validator = validator.New()
