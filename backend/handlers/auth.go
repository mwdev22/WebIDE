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
	ctr.r.Get("/login", ctr.handleGitHubLogin)
	ctr.r.Get("/callback", ctr.handleGitHubCallback)
}

func (ctr *AuthController) handleGitHubLogin(c *fiber.Ctx) error {
	url := ctr.conf.AuthCodeURL(utils.OAuthStateString)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func (ctr *AuthController) handleGitHubCallback(c *fiber.Ctx) error {
	if c.Query("state") != utils.OAuthStateString {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid oauth state"})
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "no code in query"})
	}

	token, err := ctr.conf.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to exchange token"})
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create request"})
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := ctr.client.Do(req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user info"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": resp.Status})
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to parse response"})
	}

	username := data["login"].(string)
	githubID := uint64(data["id"].(float64))
	githubURL := data["url"].(string)

	if user, err := ctr.userStore.GetUserByID(githubID); err != nil {
		fmt.Printf("user not found, %v", user)
		var newUser = types.User{
			ID:        githubID,
			GithubURL: githubURL,
			Username:  username,
		}

		if err := ctr.userStore.CreateUser(&newUser); err != nil {
			return c.JSON(fiber.Map{
				"error": err,
			})
		}
	}

	jwtToken, err := createJWT(username)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create JWT"})
	}

	return c.JSON(fiber.Map{
		"username": username,
		"gitUrl":   githubURL,
		"gitID":    githubID,
		"jwt":      jwtToken,
	})
}

func createJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(utils.SecretKey))
}
