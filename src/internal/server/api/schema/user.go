package schema

import (
	"goAuth/internal/utils/pagination"
)

type User struct {
	ID          uint8  `json:"id"`
	PhoneNumber string `json:"phone_number"`
}

// UserList uses the new generic pagination
type UserList = pagination.PaginatedResponse[User]
