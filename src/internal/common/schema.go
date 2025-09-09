package common

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

type BasicResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message,omitempty"`
}

type ErrorResponse struct {
	StatusCode  int    `json:"statusCode"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	NeedRetry   bool   `json:"need_retry,omitempty"`
	RetryReason string `json:"retry_reason,omitempty"`
}

type BasicResponseData[T any] struct {
	BasicResponse
	Data T `json:"data"`
}

var BadParamsErrorResponse = ErrorResponse{
	StatusCode: http.StatusBadRequest,
	Status:     "error",
	Message:    "bad parameters",
}

var InternalServerErrorResponse = ErrorResponse{
	StatusCode: http.StatusInternalServerError,
	Status:     "error",
	Message:    "internal server error",
}

var OkBasicResponse = BasicResponse{
	StatusCode: http.StatusOK,
	Status:     "success",
}
