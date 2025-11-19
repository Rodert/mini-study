package service

import (
	"errors"
	"time"

	"github.com/javapub/mini-study/mini-study-backend/internal/dto"
	"github.com/javapub/mini-study/mini-study-backend/internal/model"
	"github.com/javapub/mini-study/mini-study-backend/internal/repository"
	"github.com/javapub/mini-study/mini-study-backend/internal/utils"
)

// TokenService generates JWT tokens for users.
type TokenService struct {
	secret     string
	issuer     string
	ttl        time.Duration
	refreshTTL time.Duration
	userRepo   *repository.UserRepository
}

// NewTokenService creates a token service with JWT config.
func NewTokenService(secret, issuer string, ttl, refreshTTL time.Duration, userRepo *repository.UserRepository) *TokenService {
	return &TokenService{secret: secret, issuer: issuer, ttl: ttl, refreshTTL: refreshTTL, userRepo: userRepo}
}

// GeneratePair issues an access and refresh token pair.
func (s *TokenService) GeneratePair(user *model.User) (*dto.TokenResponse, error) {
	access, err := utils.GenerateToken(s.secret, s.issuer, s.ttl, user.ID, user.WorkNo)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateToken(s.secret, s.issuer, s.refreshTTL, user.ID, user.WorkNo)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{AccessToken: access, RefreshToken: refresh}, nil
}

// Refresh validates a refresh token and issues a new pair.
func (s *TokenService) Refresh(refreshToken string) (*dto.TokenResponse, error) {
	claims, err := utils.ParseToken(s.secret, refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Status {
		return nil, errors.New("用户已禁用")
	}

	return s.GeneratePair(user)
}
