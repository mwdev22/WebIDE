package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/utils"
	"golang.org/x/oauth2"
)

var conf *oauth2.Config

var jwtSecret = []byte(utils.SecretKey)

func RegisterAuth(r fiber.Router) {
	conf = utils.GetGithubConfig() // initializng the config for handlers
	r.Get("/login", handleGitHubLogin)
	r.Get("/callback", handleGitHubCallback)
}

func handleGitHubLogin(c *fiber.Ctx) error {
	url := conf.AuthCodeURL(utils.OAuthStateString)
	return c.Redirect(url, http.StatusTemporaryRedirect)
}

func handleGitHubCallback(c *fiber.Ctx) error {
	if c.Query("state") != utils.OAuthStateString {
		return c.JSON(map[string]string{"error": "invalid oauth state"})
	}

	code := c.Query("code")
	if code == "" {
		return c.JSON(map[string]string{"error": "no code in query"})
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		log.Printf("NewRequest: %s", err)
		return c.Redirect("/", http.StatusTemporaryRedirect)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do: %s", err)
		return c.Redirect("/", http.StatusTemporaryRedirect)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(map[string]string{"error": resp.Status + " returned by GithubAPI"})
	}
	var buf bytes.Buffer
	var data map[string]interface{}

	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	err = json.Unmarshal([]byte(newStr), &data)
	if err != nil {
		log.Fatal(err)
	}

	jwtToken, err := createJWT(data["login"].(string))
	if err != nil {
		log.Printf("createJWT: %s", err)
		return c.JSON(map[string]string{"error": "failed to create JWT"})
	}

	return c.JSON(map[string]string{
		"username":    data["login"].(string),
		"profile_url": data["url"].(string),
		"jwt":         jwtToken,
	})
}

func handleNewUser(c *fiber.Ctx) error {
	return nil
}

func createJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
