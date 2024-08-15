package handlers

import (
	"fmt"
	"path/filepath"
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
	// REPO ENDPOINTS
	ctr.r.Get("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleGetRepo)))
	ctr.r.Get("/user_repos/:user_id", ErrMiddleware(AuthMiddleware(ctr.handleGetUserRepos)))
	ctr.r.Get("/repo_files/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleGetRepoFiles)))

	ctr.r.Post("/new_repo", ErrMiddleware(AuthMiddleware((ctr.handleNewRepo))))

	ctr.r.Put("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleUpdateRepo)))

	// FILE ENDPOINTS
	ctr.r.Get("/file/:file_id", ErrMiddleware(AuthMiddleware(ctr.handleGetFile)))

	ctr.r.Post("/new_file", ErrMiddleware(AuthMiddleware(ctr.handleNewFile)))
	ctr.r.Post("/run_code", ErrMiddleware(AuthMiddleware(ctr.handleRunFileCode)))

	ctr.r.Put("/file/:file_id", ErrMiddleware(AuthMiddleware(ctr.handleUpdateFile)))

}

func (ctr *ProjectController) handleNewRepo(c *fiber.Ctx) error {

	var repo types.RepoPayload

	if err := c.BodyParser(&repo); err != nil {
		return InvalidJSON()
	}

	validationErrors := types.ValidateStruct(repo)
	if len(validationErrors) > 0 {
		ValidationError(validationErrors)
	}

	newRepoID, err := ctr.repoStore.CreateRepo(repo)
	if err != nil {
		return SQLError(err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      newRepoID,
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
	return c.JSON(repo)
}

func (ctr *ProjectController) handleUpdateRepo(c *fiber.Ctx) error {

	var updatedRepo types.UpdateRepoPayload

	if err := c.BodyParser(&updatedRepo); err != nil {
		return InvalidJSON()
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

	// save the updated repository to the database
	if err := ctr.repoStore.UpdateRepo(repo); err != nil {
		return SQLError(err)
	}

	// return the updated repository
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

func (ctr *ProjectController) handleNewFile(c *fiber.Ctx) error {
	var file types.FilePayload
	if err := c.BodyParser(&file); err != nil {
		return InvalidJSON()
	}

	validationErrors := types.ValidateStruct(file)
	if len(validationErrors) > 0 {
		ValidationError(validationErrors)
	}

	_, err := ctr.repoStore.GetRepoByID(file.RepositoryID)
	if err != nil {
		return NotFound(file.RepositoryID, "Repo")
	}

	extension := filepath.Ext(file.Name)
	file.Extension = extension
	runCmd := utils.GetRunCmd(extension)
	if runCmd != "" {
		file.RunCmd = runCmd + " " + file.Name // for example g++ main.cpp
	}

	newFileID, err := ctr.fileStore.CreateFile(file)
	if err != nil {
		return SQLError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      newFileID,
		"message": "File created successfully!",
	})
}

func (ctr *ProjectController) handleGetFile(c *fiber.Ctx) error {
	fileQ := c.Params("file_id")
	fileID, err := strconv.Atoi(fileQ)
	if err != nil {
		return BadQueryParameter("file_id")
	}
	file, err := ctr.fileStore.GetFileByID(fileID)
	if err != nil {
		return NotFound(fileID, "File")
	}

	return c.JSON(file)

}

func (ctr *ProjectController) handleUpdateFile(c *fiber.Ctx) error {
	var updatedFile types.UpdateFilePayload

	if err := c.BodyParser(&updatedFile); err != nil {
		return InvalidJSON()
	}

	loggedUserID, ok := c.Locals("userID").(uint)
	if !ok {
		return Unauthorized("not logged in")
	}

	fileID, err := strconv.Atoi(c.Params("file_id"))
	if err != nil {
		return BadQueryParameter("file_id")
	}

	file, err := ctr.fileStore.GetFileByID(fileID)
	if err != nil {
		return NotFound(fileID, "File")
	}

	repo, err := ctr.repoStore.GetRepoByID(file.RepositoryID)
	if err != nil {
		return NotFound(int(repo.ID), "File")
	}

	if loggedUserID != repo.UserID {
		return Unauthorized(fmt.Sprintf("User with ID %v is not the owner of repo with ID %v", loggedUserID, fileID))
	}

	if updatedFile.Name != "" && updatedFile.Name != file.Name {
		extension := filepath.Ext(file.Name)
		file.Extension = extension
		file.Name = updatedFile.Name

	}
	if updatedFile.Content != "" && updatedFile.Content != file.Content {
		file.Content = updatedFile.Content
	}

	return c.JSON(fiber.Map{
		"message": "File updated successfully",
		"file":    file,
	})
}

func (ctr *ProjectController) handleGetRepoFiles(c *fiber.Ctx) error {
	repoQ := c.Params("repo_id")
	repoID, err := strconv.Atoi(repoQ)
	if err != nil {
		return BadQueryParameter("repo_id")
	}
	repo, err := ctr.repoStore.GetRepoByID(repoID)
	if err != nil {
		return SQLError(err)
	}
	files, err := ctr.fileStore.GetFilesByRepoID(repo.ID)
	if err != nil {
		SQLError(err)
	}

	return c.JSON(fiber.Map{
		"repo_id": repoID,
		"files":   files,
	})
}

func (ctr *ProjectController) handleRunFileCode(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{})
}
