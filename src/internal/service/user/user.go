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

func (s *service) GetUsers(page, pageSize int, baseURL string, phoneNumber *string) *schema.UserList {
	page, pageSize = paginator.ValidatePagination(page, pageSize)

	query := s.db.Model(&model.User{})
	if phoneNumber != nil {
		query = query.Where("phone_number LIKE ?", "%"+*phoneNumber+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		s.logger.Error("failed to count users", zap.Error(err))
		return nil
	}

	var users []model.User
	offset := paginator.CalculateOffset(page, pageSize)
	if err := query.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		s.logger.Error("failed to list users", zap.Error(err))
		return nil
	}

	var schemaUsers []schema.User
	for _, u := range users {
		schemaUsers = append(schemaUsers, schema.User{
			ID:          u.ID,
			PhoneNumber: u.PhoneNumber,
		})
	}

	pagination := paginator.NewPagination(page, pageSize, total, baseURL)

	return paginator.NewPaginatedResponse(schemaUsers, pagination)
}
