package schema

import "github.com/go-playground/validator/v10"

type OTPRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,regex=^09[0-9]{9}$"`
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required,regex=^09[0-9]{9}$"`
	OTPCode     string `json:"otp" validate:"required,numeric"`
}

func PhoneNumberValidator(fl validator.FieldLevel) bool {
	return fl.Field().String() != "invalid"
}
