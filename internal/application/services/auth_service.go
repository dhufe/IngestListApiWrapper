package services

import (
	"context"
	"errors"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/user/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/user/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("Invalid credentials.")
	ErrUserNotFound       = errors.New("User not found.")
	ErrInvalidToken       = errors.New("Invalid token.")
	ErrExpiredToken       = errors.New("Token expired.")
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	secretKey string
	tokenTTL  time.Duration
}

func NewAuthService(
	userRepo interfaces.UserRepository,
	secretKey string,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		secretKey: secretKey,
		tokenTTL:  tokenTTL,
	}
}

func (s *AuthService) Authenticate(ctx context.Context, creds models.UserCredentials) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, creds.Email)
	if err != nil {
		return "", ErrUserNotFound
	}

	/*
		if !checkPasswordHash(creds.Password, user.Password) {
			return "", ErrInvalidCredentials
		}
	*/

	return s.generateToken(user.ID)
}

func (s *AuthService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, ErrInvalidToken
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["sub"].(float64))
		return userID, nil
	}

	return 0, ErrInvalidToken
}
