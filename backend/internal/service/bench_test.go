package service

import (
	"testing"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
)

func BenchmarkHashPassword(b *testing.B) {
	password := "benchmark-password-123"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkAuthService_GenerateToken(b *testing.B) {
	cfg := &config.JWTConfig{Secret: "bench-secret", ExpireHours: 24, RefreshExpireHours: 168}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Phone: "13800138000", Name: "Bench User", Role: model.RoleOperator}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GenerateToken(user)
	}
}
