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

	ctr.r.Post("/repo/new_repo", ErrMiddleware(AuthMiddleware((ctr.handleNewRepo))))
	ctr.r.Put("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleUpdateRepo)))
	ctr.r.Delete("/repo/:repo_id", ErrMiddleware(AuthMiddleware(ctr.handleDeleteRepo)))

	// FILE ENDPOINTS
	ctr.r.Get("/file/:file_id", ErrMiddleware(AuthMiddleware(ctr.handleGetFile)))

	ctr.r.Post("/file/new_file", ErrMiddleware(AuthMiddleware(ctr.handleNewFile)))
	ctr.r.Post("/run_code/:file_id", ErrMiddleware(AuthMiddleware(ctr.handleRunFileCode)))
	ctr.r.Put("/file/:file_id", ErrMiddleware(AuthMiddleware(ctr.handleUpdateFile)))

}

// 		REPOS

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

	if loggedUserID != repo.OwnerID {
		return Unauthorized(fmt.Sprintf("User with ID %v is not the owner of repo with ID %v", loggedUserID, repoID))
	}

	if updatedRepo.Private && updatedRepo.Private != repo.Private {
		repo.Private = updatedRepo.Private
	}
	if updatedRepo.Name != "" && updatedRepo.Name != repo.Name {
		repo.Name = updatedRepo.Name
	}
	if updatedRepo.Readme != "" && updatedRepo.Readme != repo.Readme {
		repo.Readme = updatedRepo.Readme
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

func (ctr *ProjectController) handleDeleteRepo(c *fiber.Ctx) error {
	repoQ := c.Params("repo_id")
	repoID, err := strconv.Atoi(repoQ)
	if err != nil {
		return BadQueryParameter("repo_id")
	}
	// err := ctr.repoStore.DeleteRepoByID(repoID)
	return c.JSON(fiber.Map{"id": repoID})
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
		if !repo.Private || repo.OwnerID == loggedInUserID {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	return c.JSON(fiber.Map{
		"user_id": userID,
		"repos":   filteredRepos,
	})
}

// 		FILES

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
	file.Extension = extension[1:] // removing the .
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

	if loggedUserID != repo.OwnerID {
		return Unauthorized(fmt.Sprintf("User with ID %v is not the owner of repo with ID %v", loggedUserID, fileID))
	}

	if updatedFile.Name != "" && updatedFile.Name != file.Name {
		extension := filepath.Ext(file.Name)
		file.Extension = extension[1:]
		file.Name = updatedFile.Name

	}
	if updatedFile.Content != "" && updatedFile.Content != file.Content {
		file.Content = updatedFile.Content
	}

	if err = ctr.fileStore.UpdateFile(file); err != nil {
		return SQLError(err)
	}

	return c.JSON(fiber.Map{
		"message": "File updated successfully",
		"file":    file,
	})
}

// 		COMBINED

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

	fileID, err := strconv.Atoi(c.Params("file_id"))
	if err != nil {
		return BadQueryParameter("file_id")
	}

	file, err := ctr.fileStore.GetFileByID(fileID)
	if err != nil {
		return NotFound(fileID, "file")
	}

	fileOutput, err := utils.RunCode(file)
	if err != nil {
		return NewApiError(fiber.StatusBadRequest, err)
	}

	return c.JSON(fiber.Map{
		"file_id": file.ID,
		"output":  fileOutput,
	})
}
