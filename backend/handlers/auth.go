package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/storage"
	"github.com/mwdev22/WebIDE/backend/types"
	"github.com/mwdev22/WebIDE/backend/utils"
	"golang.org/x/oauth2"
)

var jwtSecret []byte

type AuthController struct {
	r         fiber.Router
	userStore *storage.UserStore
	conf      *oauth2.Config
	client    *http.Client
}

func NewAuthController(r fiber.Router, userStore *storage.UserStore) *AuthController {
	return &AuthController{
		r:         r,
		userStore: userStore,
		conf:      utils.GetGithubConfig(),
		client:    &http.Client{},
	}
}

func (ctr *AuthController) RegisterRoutes() {
	jwtSecret = utils.SecretKey
	ctr.r.Get("/login", ErrMiddleware(ctr.handleGitHubLogin))
	ctr.r.Get("/callback", ErrMiddleware(ctr.handleGitHubCallback))
}

func (ctr *AuthController) handleGitHubLogin(c *fiber.Ctx) error {
	url := ctr.conf.AuthCodeURL(utils.OAuthStateString)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (ctr *AuthController) handleGitHubCallback(c *fiber.Ctx) error {
	if c.Query("state") != utils.OAuthStateString {
		return ExternalServiceErr(fmt.Errorf("bad state string"))
	}

	code := c.Query("code")
	if code == "" {
		return ExternalServiceErr(fmt.Errorf("no code in query"))
	}

	token, err := ctr.conf.Exchange(context.Background(), code)
	if err != nil {
		return ExternalServiceErr(err)
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return ExternalServiceErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := ctr.client.Do(req)
	if err != nil {
		return ExternalServiceErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ExternalServiceErr(err)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		InvalidJSON()
	}

	username := data["login"].(string)
	githubID := uint(data["id"].(float64))
	githubURL := data["url"].(string)

	if user, err := ctr.userStore.GetUserByID(githubID); err != nil {
		fmt.Printf("user not found, %v", user)
		var newUser = types.User{
			ID:        githubID,
			GithubURL: githubURL,
			Username:  username,
		}

		if err := ctr.userStore.CreateUser(&newUser); err != nil {
			return SQLError(err)
		}
	}

	jwtToken, err := createJWT(githubID)
	if err != nil {
		return NewApiError(fiber.StatusBadRequest, err)
	}

	return c.JSON(fiber.Map{
		"username": username,
		"gitUrl":   githubURL,
		"gitID":    githubID,
		"token":    jwtToken,
	})
}

func createJWT(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.SecretKey))
}
