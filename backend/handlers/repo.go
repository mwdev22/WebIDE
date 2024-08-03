package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/storage"
)

type RepoController struct {
	r         fiber.Router
	userStore *storage.UserStore
	repoStore *storage.RepoStore
	fileStore *storage.FileStore
}

func NewRepoController(r fiber.Router, userStore *storage.UserStore, repoStore *storage.RepoStore, fileStore *storage.FileStore) *RepoController {
	return &RepoController{
		r:         r,
		userStore: userStore,
		repoStore: repoStore,
		fileStore: fileStore,
	}
}

func (ctr *RepoController) RegisterRoutes() {
	ctr.r.Get("/user_repos/<user_id>", AuthMiddleware(ErrMiddleware(ctr.handleGetUserRepos)))
}

func (ctr *RepoController) handleGetUserRepos(c *fiber.Ctx) error {
	userID := c.Query("user_id").(int)
	userRepos, err := ctr.repoStore.GetRepoByUserID()

}
