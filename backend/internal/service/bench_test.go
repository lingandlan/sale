package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"marketplace/backend/internal/config"
	"marketplace/backend/internal/model"
)

func BenchmarkUserService_GetByID(b *testing.B) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	user := &model.User{ID: 1, Username: "testuser"}
	mockRepo.On("GetByID", ctx, int64(1)).Return(user, nil).Maybe()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.GetByID(ctx, 1)
	}
}

func BenchmarkUserService_Create(b *testing.B) {
	mockRepo := new(MockUserRepo)
	svc := NewUserService(mockRepo)
	ctx := context.Background()

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(int64(1), nil).Maybe()

	req := &model.CreateUserRequest{
		Username: "newuser",
		Password: "password123",
		Email:    "new@test.com",
		Nickname: "New User",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.Create(ctx, req)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	password := "password123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Username: "testuser", Role: model.RoleUser}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.GenerateToken(user)
	}
}

func BenchmarkParseToken(b *testing.B) {
	cfg := &config.JWTConfig{
		Secret:             "test-secret",
		ExpireHours:        1,
		RefreshExpireHours: 24,
	}
	svc := NewAuthService(cfg, nil, nil)
	user := &model.User{ID: 1, Username: "testuser", Role: model.RoleUser}

	token, _ := svc.GenerateToken(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.ParseToken(token)
	}
}

func BenchmarkGenerateTokenWithRealRSA(b *testing.B) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		b.Fatal(err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour),
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		b.Fatal(err)
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	_ = certPEM

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
	}
}
