package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/handlers"
	"github.com/mwdev22/WebIDE/backend/storage"
	"gorm.io/gorm"
)

type Server struct {
	addr string
	db   *gorm.DB
}

func NewServer(addr string, db *gorm.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {

	// grouping
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// storages
	userStore := storage.NewUserStore(s.db)
	// repoStore := storage.NewRepoStore(s.db)
	// fileStore := storage.NewFileStore(s.db)

	auth := handlers.NewAuthController(v1.Group("/auth"), userStore)
	auth.RegisterRoutes()

	return app.Listen(s.addr)
}
