package api

import (
	"errors"
	"fmt"
	"goAuth/internal/common"
	"goAuth/internal/server/api/schema"
	"goAuth/internal/utils/ratelimit"
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

// RequestOTP godoc
//
//	@Summary		Request OTP
//	@Description	Requests an OTP to be sent to the given phone number.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			OTPRequest	body		schema.OTPRequest		true	"Phone number for OTP"
//	@Success		200			{object}	common.BasicResponse	"OTP sent successfully"
//	@Failure		400			{object}	common.ErrorResponse	"Invalid request body"
//	@Failure		429			{object}	common.ErrorResponse	"Too many OTP requests"
//	@Failure		500			{object}	common.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth/request [post]
func (h *LoginHandler) RequestOTP(c *fiber.Ctx) error {
	req := new(schema.OTPRequest)
	common.Validate.RegisterValidation("regex", schema.PhoneNumberValidator)

	if errParse, errValidate := c.BodyParser(req), common.Validate.Struct(req); errParse != nil || errValidate != nil {
		h.logger.Debug("req body is not valid", zap.Any("req", req), zap.Error(errors.Join(errParse, errValidate)))
		return c.Status(http.StatusBadRequest).JSON(common.BadParamsErrorResponse)
	}

	limited, retryAfter := ratelimit.RateLimit("otp:"+req.PhoneNumber, 3, 10*60)
	if limited {
		return c.Status(http.StatusTooManyRequests).JSON(common.ErrorResponse{
			StatusCode: http.StatusTooManyRequests,
			Status:     "error",
			Message:    fmt.Sprintf("Too many OTP requests. Please try again after %d minutes.", (retryAfter+59)/60),
		})
	}

	if err := h.service.OTPRequest(req.PhoneNumber); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
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

// VerifyOTP godoc
//
//	@Summary		Verify OTP
//	@Description	Verifies the OTP code for the given phone number and returns an access token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			LoginRequest	body		schema.LoginRequest		true	"Phone number and OTP code"
//	@Success		200				{object}	common.BasicResponse	"OTP verified successfully"
//	@Failure		400				{object}	common.ErrorResponse	"Invalid request body"
//	@Failure		401				{object}	common.ErrorResponse	"Incorrect OTP code or verification failed"
//	@Failure		404				{object}	common.ErrorResponse	"OTP not found or expired"
//	@Failure		500				{object}	common.ErrorResponse	"Internal server error"
//	@Router			/api/v1/auth/verify [post]
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
			Status:     "Ok",
		},
		Data: map[string]string{"access_token": accessToken},
	})
}
