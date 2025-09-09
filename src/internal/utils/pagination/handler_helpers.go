package pagination

import (
	"github.com/gofiber/fiber/v2"
)

// ParsePaginationFromQuery extracts pagination parameters from Fiber context query parameters
func ParsePaginationFromQuery(c *fiber.Ctx) PaginationParams {
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("page_size", "10")

	return ParsePaginationParams(pageStr, pageSizeStr)
}

// GetBaseURL constructs the base URL for pagination links
func GetBaseURL(c *fiber.Ctx) string {
	scheme := "http"
	if c.Secure() {
		scheme = "https"
	}

	host := c.Hostname()
	path := c.Path()

	return scheme + "://" + host + path
}
