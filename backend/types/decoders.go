package types

type RegisterUserPayload struct {
	Username  string `json:"username" validate:"required"`
	GithubID  string `json:"gitId" validate:"required"`
	GithubURL string `json:"gitUrl" validate:"required"`
}
