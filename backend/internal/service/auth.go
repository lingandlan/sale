package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
	"marketplace/backend/internal/repository"
	apperrors "marketplace/backend/pkg/errors"
)

// Claims JWT Claims
type Claims struct {
	UserID int64  `json:"user_id"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService 认证服务
type AuthService struct {
	jwtSecret          []byte
	expireHours        time.Duration
	refreshExpireHours time.Duration
	redis              *redis.Client
	userRepo           repository.UserRepoInterface
}

// NewAuthService 创建 AuthService
func NewAuthService(cfg *config.JWTConfig, redisClient *redis.Client, userRepo repository.UserRepoInterface) *AuthService {
	return &AuthService{
		jwtSecret:          []byte(cfg.Secret),
		expireHours:        time.Duration(cfg.ExpireHours) * time.Hour,
		refreshExpireHours: time.Duration(cfg.RefreshExpireHours) * time.Hour,
		redis:              redisClient,
		userRepo:           userRepo,
	}
}

// GenerateToken 生成 Access Token
func (s *AuthService) GenerateToken(user *model.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Phone:  user.Phone,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expireHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// GenerateRefreshToken 生成 Refresh Token
func (s *AuthService) GenerateRefreshToken(user *model.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Phone:  user.Phone,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshExpireHours)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ParseToken 解析 Token
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.ErrTokenInvalid
}

// StoreRefreshToken 存储 Refresh Token 到 Redis
func (s *AuthService) StoreRefreshToken(ctx context.Context, userID int64, token string) error {
	return repository.StoreRefreshToken(ctx, s.redis, userID, token)
}

// ValidateRefreshToken 验证 Refresh Token
func (s *AuthService) ValidateRefreshToken(ctx context.Context, userID int64, token string) (bool, error) {
	return repository.ValidateRefreshToken(ctx, s.redis, userID, token)
}

// RevokeRefreshToken 撤销 Refresh Token
func (s *AuthService) RevokeRefreshToken(ctx context.Context, userID int64) error {
	return repository.RevokeRefreshToken(ctx, s.redis, userID)
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, phone, password string) (*model.LoginResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, apperrors.ErrNotFound
	}

	// 检查用户状态
	if user.Status != model.UserStatusNormal {
		return nil, apperrors.ErrUserDisabled
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, apperrors.ErrPasswordIncorrect
	}

	// 生成 Token
	accessToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate access token failed: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token failed: %w", err)
	}

	// 存储 Refresh Token
	if err := s.StoreRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("store refresh token failed: %w", err)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.expireHours.Seconds()),
	}, nil
}

// RefreshToken 刷新 Token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// 解析 Token
	claims, err := s.ParseToken(refreshToken)
	if err != nil {
		return "", "", apperrors.ErrTokenInvalid
	}

	// 验证 Redis 中的 Token
	valid, err := s.ValidateRefreshToken(ctx, claims.UserID, refreshToken)
	if err != nil {
		return "", "", err
	}
	if !valid {
		return "", "", apperrors.ErrTokenInvalid
	}

	// 获取用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", apperrors.ErrNotFound
	}

	// 生成新 Token
	newAccessToken, err := s.GenerateToken(user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	// 更新 Refresh Token
	if err := s.StoreRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// Logout 用户登出
func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	return s.RevokeRefreshToken(ctx, userID)
}

// HashPassword 密码哈希
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
