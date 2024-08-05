package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/storage"
	"github.com/mwdev22/WebIDE/backend/types"
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
	ctr.r.Post("/new_repo", ErrMiddleware(AuthMiddleware((ctr.handleNewRepo))))
	ctr.r.Get("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleGetRepo)))
	ctr.r.Get("/user_repos/:user_id", ErrMiddleware(AuthMiddleware(ctr.handleGetUserRepos)))

}

func (ctr *RepoController) handleNewRepo(c *fiber.Ctx) error {

	var repo types.Repo

	if err := c.BodyParser(&repo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body.",
			"error":   err.Error(),
		})
	}

	validationErrors := types.ValidateStruct(repo)
	if len(validationErrors) > 0 {
		ValidationError(validationErrors)
	}

	newRepoID, err := ctr.repoStore.CreateRepo(repo)
	if err != nil {
		return SQLError(err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"repo_id": newRepoID,
		"message": "Repository created successfully.",
	})
}

func (ctr *RepoController) handleGetRepo(c *fiber.Ctx) error {
	repoQ := c.Params("repo_id")
	repoID, err := strconv.Atoi(repoQ)
	if err != nil {
		return BadQueryParameter("repo_id")
	}
	repo, err := ctr.repoStore.GetRepoByID(repoID)
	if err != nil {
		return NotFound(repoID, "Repository")
	}
	return c.Status(fiber.StatusOK).JSON(repo)
}

func (ctr *RepoController) handleGetUserRepos(c *fiber.Ctx) error {
	userQ := c.Params("user_id")
	userID, err := strconv.Atoi(userQ)
	if err != nil {
		BadQueryParameter("user_id")
	}
	userRepos, err := ctr.repoStore.GetReposByUserID(userID)
	if err != nil {
		return SQLError(err)
	}

	loggedInUserID, ok := c.Locals("userID").(uint)
	if !ok {
		return Unauthorized("not logged in")
	}
	// return only private repos or all if the logged user is an owner
	filteredRepos := []*storage.Repository{}
	for _, repo := range userRepos {
		if !repo.Private || repo.UserID == loggedInUserID {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"repos":   filteredRepos,
	})
}
