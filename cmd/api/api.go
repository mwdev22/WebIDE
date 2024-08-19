package api

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/websocket/v2"
	"github.com/mwdev22/WebIDE/cmd/handlers"
	"github.com/mwdev22/WebIDE/cmd/storage"
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
	app.Use(logger.New(logger.Config{
		Output:     os.Stdout,
		Format:     "${time} ${ip} ${status} - ${method} ${path}\n",
		TimeFormat: "31-02-2006 15:04:05",
		TimeZone:   "Europe/Warsaw",
	}))

	// storages
	userStore := storage.NewUserStore(s.db)
	repoStore := storage.NewRepoStore(s.db)
	fileStore := storage.NewFileStore(s.db)

	auth := handlers.NewAuthController(v1.Group("/auth"), userStore)
	project := handlers.NewProjectController(v1, userStore, repoStore, fileStore)

	auth.RegisterRoutes()
	project.RegisterRoutes()

	// websocket for colaborating with other users
	app.Get("/ws/:fileId", websocket.New(func(c *websocket.Conn) {
		handlers.HandleWebSocketConnection(c)
	}))

	return app.Listen(s.addr)
}
