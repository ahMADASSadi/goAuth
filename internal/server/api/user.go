package api

import (
	"goAuth/internal/utils/pagination"
	"goAuth/internal/server/api/schema"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserService interface {
	GetUser(id uint8) *schema.User
	GetUsers(page, pageSize int, baseURL string) *schema.UserList
}

type UserHandler struct {
	logger  *zap.Logger
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		logger:  zap.L(),
		service: service,
	}
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "user id is required",
		})
	}
	idUint, err := strconv.ParseUint(idStr, 10, 8)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user id",
		})
	}
	user := h.service.GetUser(uint8(idUint))
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "user not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse pagination parameters from query
	params := pagination.ParsePaginationFromQuery(c)

	// Get base URL for pagination links
	baseURL := pagination.GetBaseURL(c)

	// Get paginated users
	users := h.service.GetUsers(params.Page, params.PageSize, baseURL)
	if users == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to get users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
