package api

import (
	"goAuth/internal/server/api/schema"
	"goAuth/internal/utils/pagination"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserService interface {
	GetUser(id uint8) *schema.User
	GetUsers(page, pageSize int, baseURL string, phoneNumber *string) *schema.UserList
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

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Retrieves a user by their unique ID.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	schema.User "returns the user details including the id and phone number"
//	@Failure		400	{object}	common.ErrorResponse
//	@Failure		404	{object}	common.ErrorResponse
//	@Router			/api/v1/users/{id} [get]
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

// GetUsers godoc
//
//	@Summary		List users
//	@Description	Retrieves a paginated list of users, with optional phone number search.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			page			query	int		false	"Page number"				default(1)
//	@Param			page_size		query	int		false	"Number of users per page"	default(10)
//	@Param			phone_number	query	string	false	"Filter by phone number"
//	@Success		200				{array}	schema.User
//	@Router			/api/v1/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	var phoneNumber *string
	// Parse pagination parameters from query
	params := pagination.ParsePaginationFromQuery(c)
	phoneNumberQuery := c.Query("phone_number")
	if phoneNumberQuery != "" {
		phoneNumber = &phoneNumberQuery
	}

	baseURL := pagination.GetBaseURL(c)

	// Get paginated users
	users := h.service.GetUsers(params.Page, params.PageSize, baseURL, phoneNumber)
	if users == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to get users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
