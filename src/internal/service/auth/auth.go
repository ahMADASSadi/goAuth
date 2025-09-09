package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"goAuth/internal/common"
	"goAuth/internal/database/model"
	inmemory "goAuth/internal/service/in-memory"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	db     *gorm.DB
	logger *zap.Logger
	inMemo *inmemory.InMemoryStore
}

func NewAuthenticationService(db *gorm.DB, inMemo *inmemory.InMemoryStore) *service {
	return &service{
		db:     db,
		logger: zap.L(),
		inMemo: inMemo,
	}
}

func (s *service) OTPRequest(phoneNumber string) error {
	max := big.NewInt(1000000) // 0 to 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		s.logger.Error("failed to generate OTP", zap.Error(err))
		return fmt.Errorf("failed to generate OTP: %w", err)
	}
	otpCode := fmt.Sprintf("%06s", n.String())

	s.inMemo.Set(phoneNumber, otpCode, time.Minute*2) // 2 minutes

	s.logger.Debug("OTP generated", zap.String("phoneNumber", phoneNumber), zap.String("otp", otpCode))
	return nil
}

func (s *service) OTPVerify(phoneNumber, otpCode string) (bool, error) {
	registeredOTP, ok := s.inMemo.Get(phoneNumber)
	if !ok {
		return false, common.ErrGetOTP
	}

	otpStr, ok := registeredOTP.(string)
	if !ok {
		s.logger.Error("stored OTP is not a string", zap.Any("registeredOTP", registeredOTP))
		return false, common.ErrInvalidOTP
	}

	if otpStr != otpCode {
		return false, common.ErrCompareOTP
	}
	return true, nil
}

func (s *service) RegisterUser(phoneNumber string) (created bool, err error) {
	var user model.User

	dbErr := s.db.Where("phone_number = ?", phoneNumber).First(&user).Error

	if dbErr == nil {
		return false, nil
	}

	if !errors.Is(dbErr, gorm.ErrRecordNotFound) {
		s.logger.Error("database error checking for user", zap.Error(dbErr))
		return false, dbErr
	}

	newUser := &model.User{
		PhoneNumber: phoneNumber,
	}
	if createErr := s.db.Create(newUser).Error; createErr != nil {
		s.logger.Error("failed to create user", zap.Error(createErr), zap.String("phoneNumber", phoneNumber))
		return false, createErr
	}
	return true, nil

}

func (s *service) GenerateToken(phoneNumber string) (accessToken string, err error) {
	expiryStr := os.Getenv("ACCESS_EXPIRY")
	secretStr := os.Getenv("SECRET_KEY")
	expiryDuration, err := time.ParseDuration(expiryStr)

	if err != nil {
		s.logger.Error("invalid accessExpiry duration", zap.String("accessExpiry", expiryStr), zap.Error(err))
		return "", fmt.Errorf("invalid accessExpiry duration: %w", err)
	}

	iat := time.Now().Unix()

	hasher := sha256.New()
	hasher.Write([]byte(phoneNumber + fmt.Sprint(iat)))
	userHash := fmt.Sprintf("%x", hasher.Sum(nil))

	claims := jwt.MapClaims{
		"user": userHash,
		"exp":  time.Now().Add(expiryDuration).Unix(),
		"iat":  iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretStr))
	if err != nil {
		s.logger.Error("failed to sign JWT token", zap.Error(err))
		return "", err
	}

	return signedToken, nil
}
