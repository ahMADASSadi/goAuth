package model

import "time"

type User struct {
	ID          uint8 `gorm:"primarykey"`
	CreatedAt   time.Time
	PhoneNumber string `gorm:"unique;not null" validate:"required,regexp=^09[0-9]{9}$"`
}
