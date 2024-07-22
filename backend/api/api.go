package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mwdev22/WebIDE/backend/handlers"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Run() error {
	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	handlers.RegisterAuth(auth)

	return app.Listen(s.addr)
}
