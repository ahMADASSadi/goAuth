package pagination

import (
	"net/url"
	"strconv"
)

// Pagination represents pagination metadata
type Pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	NextPage   string `json:"next_page,omitempty"`
	PrevPage   string `json:"prev_page,omitempty"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}

// PaginatedResponse represents a generic paginated response
type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page     int
	PageSize int
}

// DefaultPagination returns default pagination parameters
func DefaultPagination() PaginationParams {
	return PaginationParams{
		Page:     1,
		PageSize: 10,
	}
}

// ParsePaginationParams parses pagination parameters from query values
func ParsePaginationParams(pageStr, pageSizeStr string) PaginationParams {
	page := 1
	pageSize := 10

	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
		pageSize = ps
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// NewPagination creates a new pagination object with calculated metadata
func NewPagination(page, pageSize int, total int64, baseURL string) Pagination {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	pagination := Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}

	// Generate next page URL
	if pagination.HasNext {
		pagination.NextPage = buildPageURL(baseURL, page+1, pageSize)
	}

	// Generate previous page URL
	if pagination.HasPrev {
		pagination.PrevPage = buildPageURL(baseURL, page-1, pageSize)
	}

	return pagination
}

// buildPageURL constructs a URL with pagination parameters
func buildPageURL(baseURL string, page, pageSize int) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("page_size", strconv.Itoa(pageSize))
	u.RawQuery = q.Encode()

	return u.String()
}

// CalculateOffset calculates the database offset for pagination
func CalculateOffset(page, pageSize int) int {
	if page <= 0 {
		page = 1
	}
	return (page - 1) * pageSize
}

// ValidatePagination validates and normalizes pagination parameters
func ValidatePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, pagination Pagination) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{
		Data:       data,
		Pagination: pagination,
	}
}
