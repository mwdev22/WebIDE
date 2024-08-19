package utils

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var OAuthStateString = "random"

func GetGithubConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}
