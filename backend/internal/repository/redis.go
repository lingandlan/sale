package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"marketplace/backend/internal/config"
)

// NewRedis 创建 Redis 客户端
func NewRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 100,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis failed: %w", err)
	}

	return client, nil
}

// Redis 常量
const (
	RefreshTokenPrefix = "refresh_token:"
	RefreshTokenTTL    = 7 * 24 * time.Hour // 7 天
)

// StoreRefreshToken 存储 Refresh Token
func StoreRefreshToken(ctx context.Context, rdb *redis.Client, userID int64, token string) error {
	key := fmt.Sprintf("%s%d", RefreshTokenPrefix, userID)
	return rdb.Set(ctx, key, token, RefreshTokenTTL).Err()
}

// ValidateRefreshToken 验证 Refresh Token
func ValidateRefreshToken(ctx context.Context, rdb *redis.Client, userID int64, token string) (bool, error) {
	key := fmt.Sprintf("%s%d", RefreshTokenPrefix, userID)
	stored, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return stored == token, nil
}

// RevokeRefreshToken 撤销 Refresh Token
func RevokeRefreshToken(ctx context.Context, rdb *redis.Client, userID int64) error {
	key := fmt.Sprintf("%s%d", RefreshTokenPrefix, userID)
	return rdb.Del(ctx, key).Err()
}
