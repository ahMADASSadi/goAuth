package user

import (
	"goAuth/internal/database/model"
	"goAuth/internal/server/api/schema"
	paginator "goAuth/internal/utils/pagination"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewUserService(db *gorm.DB) *service {
	return &service{
		db:     db,
		logger: zap.L(),
	}
}

func (s *service) GetUser(id uint8) *schema.User {
	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		s.logger.Error("failed to get user", zap.Error(err))
		return nil
	}
	return &schema.User{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
	}
}

func (s *service) GetUsers(page, pageSize int, baseURL string) *schema.UserList {
	// Validate and normalize pagination parameters
	page, pageSize = paginator.ValidatePagination(page, pageSize)

	// Count total records
	var total int64
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		s.logger.Error("failed to count users", zap.Error(err))
		return nil
	}

	// Get paginated data
	var users []model.User
	offset := paginator.CalculateOffset(page, pageSize)
	if err := s.db.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil
	}

	// Convert to schema
	var schemaUsers []schema.User
	for _, u := range users {
		schemaUsers = append(schemaUsers, schema.User{
			ID:          u.ID,
			PhoneNumber: u.PhoneNumber,
		})
	}

	// Create pagination metadata
	pagination := paginator.NewPagination(page, pageSize, total, baseURL)

	// Return paginated response
	return paginator.NewPaginatedResponse(schemaUsers, pagination)
}
