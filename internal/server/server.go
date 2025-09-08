package server

import (
	"github.com/gofiber/fiber/v2"

	"goAuth/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "goAuth",
			AppName:      "goAuth",
		}),

		db: database.New(),
	}

	return server
}
