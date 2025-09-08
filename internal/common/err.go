package common

import "errors"

var (
	ErrGetOTP     = errors.New("fetching registered otp code faild")
	ErrCompareOTP = errors.New("wrong otp code")
	ErrInvalidOTP = errors.New("invalid otp")
)
