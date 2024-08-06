package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/storage"
	"github.com/mwdev22/WebIDE/backend/types"
	"github.com/mwdev22/WebIDE/backend/utils"
)

type ProjectController struct {
	r         fiber.Router
	userStore *storage.UserStore
	repoStore *storage.RepoStore
	fileStore *storage.FileStore
}

func NewProjectController(r fiber.Router, userStore *storage.UserStore, repoStore *storage.RepoStore, fileStore *storage.FileStore) *ProjectController {
	return &ProjectController{
		r:         r,
		userStore: userStore,
		repoStore: repoStore,
		fileStore: fileStore,
	}
}

func (ctr *ProjectController) RegisterRoutes() {
	ctr.r.Get("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleGetRepo)))
	ctr.r.Get("/user_repos/:user_id", ErrMiddleware(AuthMiddleware(ctr.handleGetUserRepos)))

	ctr.r.Post("/new_repo", ErrMiddleware(AuthMiddleware((ctr.handleNewRepo))))

	ctr.r.Patch("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleUpdateRepo)))
}

func (ctr *ProjectController) handleNewRepo(c *fiber.Ctx) error {

	var repo types.RepoPayload

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

func (ctr *ProjectController) handleGetRepo(c *fiber.Ctx) error {
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

func (ctr *ProjectController) handleUpdateRepo(c *fiber.Ctx) error {

	var updatedRepo types.UpdateRepoPayload

	if err := c.BodyParser(&updatedRepo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to parse request body.",
			"error":   err.Error(),
		})
	}

	loggedUserID, ok := c.Locals("userID").(uint)
	if !ok {
		return Unauthorized("not logged in")
	}

	repoID, err := strconv.Atoi(c.Params("repo_id"))
	if err != nil {
		return BadQueryParameter("repo_id")
	}

	repo, err := ctr.repoStore.GetRepoByID(repoID)
	if err != nil {
		return NotFound(repoID, "Repository")
	}

	if loggedUserID != repo.UserID {
		return Unauthorized(fmt.Sprintf("User with ID %v is not the owner of repo with ID %v", loggedUserID, repoID))
	}

	if err = utils.CheckAndUpdate(updatedRepo, &repo); err != nil {
		return NewApiError(fiber.StatusBadGateway, err)
	}

	// Save the updated repository to the database
	if err := ctr.repoStore.UpdateRepo(repo); err != nil {
		return SQLError(err)
	}

	// Return the updated repository
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Repository updated successfully",
		"repo":    repo,
	})
}

func (ctr *ProjectController) handleGetUserRepos(c *fiber.Ctx) error {
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
