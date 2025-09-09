package server

import (
	"github.com/gofiber/fiber/v2"

	"goAuth/internal/database"
)

type FiberServer struct {
	*fiber.App

	DB database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "goAuth",
			AppName:      "goAuth",
		}),

		DB: database.New(),
	}

	return server
}
