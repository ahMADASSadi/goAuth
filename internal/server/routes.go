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
		AllowCredentials: false, // credentials require explicit origins
		MaxAge:           300,
	}))
	setupSwagger(s.App)
	setupAuthRoutes(s.App, authService)
	setupUserRoutes(s.App, userService)

	s.App.Get("/", s.healthRoutes)

}

func (s *FiberServer) healthRoutes(c *fiber.Ctx) error {
	return c.JSON(s.DB.Health())
}

func setupSwagger(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault) // default

}

func setupAuthRoutes(app *fiber.App, service api.LoginService) {
	handler := api.NewLoginHandler(service)

	// @Summary Send OTP
	// @Description Sends a one-time password to the user's phone number.
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param data body api.RequestOTPRequest true "Phone number"
	// @Success 200 {object} api.RequestOTPResponse
	// @Failure 400 {object} api.ErrorResponse
	// @Router /send-otp [post]
	app.Post("/send-otp", handler.RequestOTP)

	// @Summary Verify OTP
	// @Description Verifies the one-time password sent to the user's phone number.
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param data body api.VerifyOTPRequest true "Phone number and OTP code"
	// @Success 200 {object} api.VerifyOTPResponse
	// @Failure 400 {object} api.ErrorResponse
	// @Router /verify-otp [post]
	app.Post("/verify-otp", handler.VerifyOTP)
}

func setupUserRoutes(app *fiber.App, service api.UserService) {

	handler := api.NewUserHandler(service)
	// Get user by ID
	app.Get("/users/:id", handler.GetUser)

	// List users with optional phone number search and pagination
	app.Get("/users", handler.GetUsers)
}
