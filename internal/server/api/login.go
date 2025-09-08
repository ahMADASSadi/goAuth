package api

import (
	"errors"
	"fmt"
	"goAuth/internal/common"
	"goAuth/internal/server/api/schema"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type LoginService interface {
	OTPRequest(phoneNumber string) error
	OTPVerify(phoneNumber string, otpCode string) (bool, error)
	RegisterUser(phoneNumber string) (created bool, err error)
	GenerateToken(phoneNumber string) (accessToken string, err error)
}

type LoginHandler struct {
	logger  *zap.Logger
	service LoginService
}

func NewLoginHandler(service LoginService) *LoginHandler {
	return &LoginHandler{
		logger:  zap.L(),
		service: service,
	}
}

func (h *LoginHandler) RequestOTP(c *fiber.Ctx) error {
	req := new(schema.OTPRequest)
	common.Validate.RegisterValidation("regex", schema.PhoneNumberValidator)

	if errParse, errValidate := c.BodyParser(req), common.Validate.Struct(req); errParse != nil || errValidate != nil {
		h.logger.Debug("req body is not valid", zap.Any("req", req), zap.Error(errors.Join(errParse, errValidate)))
		return c.Status(http.StatusBadRequest).JSON(common.BadParamsErrorResponse)
	}
	if err := h.service.OTPRequest(req.PhoneNumber); err != nil {
		c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "OTP Creation faild",
		})
	}
	return c.Status(http.StatusOK).JSON(common.BasicResponse{
		StatusCode: http.StatusOK,
		Message:    "OTP sent successfully",
	})
}

func (h *LoginHandler) VerifyOTP(c *fiber.Ctx) error {
	common.Validate.RegisterValidation("regex", schema.PhoneNumberValidator)

	req := new(schema.LoginRequest)
	if errParse, errValidate := c.BodyParser(req), common.Validate.Struct(req); errParse != nil || errValidate != nil {
		h.logger.Debug("req body is not valid", zap.Any("req", req), zap.Error(errors.Join(errParse, errValidate)))
		return c.Status(http.StatusBadRequest).JSON(common.BadParamsErrorResponse)
	}

	verified, err := h.service.OTPVerify(req.PhoneNumber, req.OTPCode)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrGetOTP):
			return c.Status(http.StatusNotFound).JSON(common.ErrorResponse{
				StatusCode: http.StatusNotFound,
				Status:     "error",
				Message:    "OTP not found or expired",
			})
		case errors.Is(err, common.ErrInvalidOTP):
			return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "Internal error verifying OTP",
			})
		case errors.Is(err, common.ErrCompareOTP):
			return c.Status(http.StatusUnauthorized).JSON(common.ErrorResponse{
				StatusCode: http.StatusUnauthorized,
				Status:     "error",
				Message:    "Incorrect OTP code",
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				Status:     "error",
				Message:    "Unknown error verifying OTP",
			})
		}
	}

	if !verified {
		return c.Status(http.StatusUnauthorized).JSON(common.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "OTP verification failed",
		})
	}

	_, err = h.service.RegisterUser(req.PhoneNumber)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    "Error creating new user",
		})
	}
	accessToken, err := h.service.GenerateToken(req.PhoneNumber)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "error",
			Message:    fmt.Sprintf("err: %v\n", err),
		})
	}

	return c.Status(http.StatusOK).JSON(common.BasicResponseData[map[string]string]{
		BasicResponse: common.BasicResponse{
			StatusCode: http.StatusOK,
			Message:    "OTP verified successfully",
		},
		Data: map[string]string{"access_token": accessToken},
	})
}
