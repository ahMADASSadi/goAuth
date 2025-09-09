package server

import (
	"goAuth/internal/server/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	_ "goAuth/docs"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) SetupRoutes(authService api.LoginService, userService api.UserService) {
	// Apply CORS middleware
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	apiV1 := s.App.Group("/api/v1")
	apiV1.Get("/", s.healthRoutes)

	// Swagger docs route
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	// Auth routes: /api/v1/auth/request, /api/v1/auth/verify
	authGroup := apiV1.Group("/auth")
	setupAuthRoutes(authGroup, authService)

	// User routes: /api/v1/users/:id, /api/v1/users
	setupUserRoutes(apiV1, userService)
}

func (s *FiberServer) healthRoutes(c *fiber.Ctx) error {
	return c.JSON(s.DB.Health())
}

func setupAuthRoutes(app fiber.Router, service api.LoginService) {
	handler := api.NewLoginHandler(service)

	// POST /api/v1/auth/request
	app.Post("/request", handler.RequestOTP)

	// POST /api/v1/auth/verify
	app.Post("/verify", handler.VerifyOTP)
}

func setupUserRoutes(app fiber.Router, service api.UserService) {
	handler := api.NewUserHandler(service)

	// GET /api/v1/users/:id
	app.Get("/users/:id", handler.GetUser)

	// GET /api/v1/users
	app.Get("/users", handler.GetUsers)
}
